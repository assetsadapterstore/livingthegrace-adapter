/*
 * Copyright 2019 The openwallet Authors
 * This file is part of the openwallet library.
 *
 * The openwallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The openwallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package livingthegrace

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/assetsadapterstore/livingthegrace-adapter/extends/api/p2p"
	"github.com/assetsadapterstore/livingthegrace-adapter/extends/crypto/rsencoding"
	"github.com/assetsadapterstore/livingthegrace-adapter/extends/transaction"
	"github.com/blocktree/go-owcrypt"
	"github.com/blocktree/openwallet/v2/openwallet"
	"github.com/shopspring/decimal"
	"strconv"
	"time"
)

type TransactionDecoder struct {
	openwallet.TransactionDecoderBase
	wm *WalletManager //钱包管理者
}

//NewTransactionDecoder 交易单解析器
func NewTransactionDecoder(wm *WalletManager) *TransactionDecoder {
	decoder := TransactionDecoder{}
	decoder.wm = wm
	return &decoder
}

//CreateRawTransaction 创建交易单
func (decoder *TransactionDecoder) CreateRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) error {

	var (
		accountID       = rawTx.Account.AccountID
		findAddrBalance *LTGAccount
	)

	//获取wallet
	addresses, err := wrapper.GetAddressList(0, -1, "AccountID", accountID) //wrapper.GetWallet().GetAddressesByAccount(rawTx.Account.AccountID)
	if err != nil {
		return err
	}

	if len(addresses) == 0 {
		return openwallet.Errorf(openwallet.ErrAccountNotAddress, "[%s] have not addresses", accountID)
	}

	var amountStr string
	for _, v := range rawTx.To {
		amountStr = v
		break
	}

	amount, _ := decimal.NewFromString(amountStr)



	fee := decoder.GetFees()
	for _, addr := range addresses {

		addrBalance, err := decoder.wm.GetWalletDetails(addr.Address)
		if err != nil {
			continue
		}

		balance, _ := decimal.NewFromString(addrBalance.Amount)
		balance = balance.Shift(-decoder.wm.Decimal())
		//余额不足查找下一个地址
		totalSend := amount.Add(fee)
		if balance.GreaterThanOrEqual(totalSend) {
			//只要找到一个合适使用的地址余额就停止遍历
			findAddrBalance = addrBalance
			break
		}
	}
	rawTx.Fees = fee.String()

	if findAddrBalance == nil {
		return openwallet.Errorf(openwallet.ErrInsufficientBalanceOfAccount, "all address's balance of account is not enough")
	}

	//最后创建交易单
	err = decoder.createRawTransaction(
		wrapper,
		rawTx,
		findAddrBalance)
	if err != nil {
		return err
	}

	return nil

}


func (decoder *TransactionDecoder) GetFees() decimal.Decimal{
	supportFee,err := decoder.wm.GetSupportFee()
	if err != nil{
		return decoder.wm.Config.FixFees
	}
	feeLong := supportFee.Standard
	feeDecimal := decimal.NewFromInt(int64(feeLong)).Shift(-decoder.wm.Decimal())
	return feeDecimal
}

//SignRawTransaction 签名交易单
func (decoder *TransactionDecoder) SignRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) error {

	if rawTx.Signatures == nil || len(rawTx.Signatures) == 0 {
		//this.wm.Log.Std.Error("len of signatures error. ")
		return fmt.Errorf("transaction signature is empty")
	}

	key, err := wrapper.HDKey()
	if err != nil {
		return err
	}

	keySignatures := rawTx.Signatures[rawTx.Account.AccountID]
	if keySignatures != nil {
		for _, keySignature := range keySignatures {

			childKey, err := key.DerivedKeyWithPath(keySignature.Address.HDPath, keySignature.EccType)
			keyBytes, err := childKey.GetPrivateKeyBytes()
			if err != nil {
				return err
			}

			msg, err := hex.DecodeString(keySignature.Message)
			if err != nil {
				return fmt.Errorf("decoder transaction hash failed, unexpected err: %v", err)
			}

			sig, _, ret := owcrypt.Signature(keyBytes, nil, msg, keySignature.EccType)
			if ret != owcrypt.SUCCESS {
				return fmt.Errorf("sign transaction hash failed, unexpected err: %v", err)
			}

			keySignature.Signature = hex.EncodeToString(sig)
		}
	}

	decoder.wm.Log.Info("transaction hash sign success")

	rawTx.Signatures[rawTx.Account.AccountID] = keySignatures

	return nil
}

//VerifyRawTransaction 验证交易单，验证交易单并返回加入签名后的交易单
func (decoder *TransactionDecoder) VerifyRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) error {

	if rawTx.Signatures == nil || len(rawTx.Signatures) == 0 {
		//this.wm.Log.Std.Error("len of signatures error. ")
		return fmt.Errorf("transaction signature is empty")
	}

	txHex, err := hex.DecodeString(rawTx.RawHex)
	if err != nil {
		return fmt.Errorf("transaction decode failed, unexpected error: %v", err)
	}

	var tx RawTransactionV2
	err = json.Unmarshal(txHex, &tx)
	if err != nil {
		return err
	}

	//支持多重签名
	for accountID, keySignatures := range rawTx.Signatures {
		decoder.wm.Log.Debug("accountID Signatures:", accountID)
		for _, keySignature := range keySignatures {

			//messsage, _ := hex.DecodeString(keySignature.Message)
			signature, _ := hex.DecodeString(keySignature.Signature)
			publicKey, _ := hex.DecodeString(keySignature.Address.PublicKey)
			msg, _ := hex.DecodeString(keySignature.Message)

			fmt.Println("msg:", keySignature.Message)

			enpub, _ := owcrypt.CURVE25519_convert_Ed_to_X(publicKey)

			////验证签名
			ret := owcrypt.Verify(enpub, nil, msg, signature, keySignature.EccType)
			if ret != owcrypt.SUCCESS {
				return fmt.Errorf("transaction verify failed")
			}

			tNewSend := RawTransactionSend{
				SenderPublicKey: tx.Sender,
				Signature:       hex.EncodeToString(signature),
				Recipient:       tx.Recipient,
				AmountNQT:       tx.Amount,
				FeeNQT:          tx.Fee,
				Deadline:        tx.DeadLine,
				Timestamp:       tx.Timestamp,
				Version:         tx.Version,
				EcBlockHeight:   tx.EcBlockHeight,
				EcBlockId:       tx.EcBlockId,
				Appendix:        tx.Appendix,
			}

			result, _ := json.Marshal(tNewSend)

			rawTx.RawHex = string(result)
			rawTx.IsCompleted = true

		}
	}

	return nil
}

//SendRawTransaction 广播交易单
func (decoder *TransactionDecoder) SubmitRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) (*openwallet.Transaction, error) {

	txID, err := decoder.wm.SendOffTx(rawTx.RawHex)
	if err != nil {
		//decoder.wm.UpdateAddressNonce(wrapper, txSigned.Sender, 0)
		return nil, err
	}

	rawTx.TxID = txID
	//交易成功，地址nonce+1并记录到缓存
	//decoder.wm.UpdateAddressNonce(wrapper, txSigned.Sender, txSigned.Nonce)

	//decoder.wm.Log.Infof("Transaction [%s] submitted to the network successfully.", txid)

	rawTx.IsSubmit = true

	decimals := decoder.wm.Decimal()

	//记录一个交易单
	tx := &openwallet.Transaction{
		From:       rawTx.TxFrom,
		To:         rawTx.TxTo,
		Amount:     rawTx.TxAmount,
		Coin:       rawTx.Coin,
		TxID:       rawTx.TxID,
		Decimal:    decimals,
		AccountID:  rawTx.Account.AccountID,
		Fees:       rawTx.Fees,
		SubmitTime: time.Now().Unix(),
		ExtParam:   rawTx.ExtParam,
	}

	tx.WxID = openwallet.GenTransactionWxID(tx)

	return tx, nil
}

//GetRawTransactionFeeRate 获取交易单的费率
func (decoder *TransactionDecoder) GetRawTransactionFeeRate() (feeRate string, unit string, err error) {
	return decoder.GetFees().String(), "TX", nil
}

//CreateSummaryRawTransaction 创建汇总交易
func (decoder *TransactionDecoder) CreateSummaryRawTransaction(wrapper openwallet.WalletDAI, sumRawTx *openwallet.SummaryRawTransaction) ([]*openwallet.RawTransaction, error) {

	var (
		rawTxArray         = make([]*openwallet.RawTransaction, 0)
		accountID          = sumRawTx.Account.AccountID
		minTransfer, _     = decimal.NewFromString(sumRawTx.MinTransfer)
		retainedBalance, _ = decimal.NewFromString(sumRawTx.RetainedBalance)
	)

	if minTransfer.Cmp(retainedBalance) < 0 {
		return nil, fmt.Errorf("mini transfer amount must be greater than address retained balance")
	}

	//获取wallet
	addresses, err := wrapper.GetAddressList(sumRawTx.AddressStartIndex, sumRawTx.AddressLimit,
		"AccountID", sumRawTx.Account.AccountID)
	if err != nil {
		return nil, err
	}

	if len(addresses) == 0 {
		return nil, fmt.Errorf("[%s] have not addresses", accountID)
	}


	for _, addr := range addresses {


		fee := decoder.GetFees()
		addrBalance, err := decoder.wm.GetWalletDetails(addr.Address)
		if err != nil {
			continue
		}

		balance, _ := decimal.NewFromString(addrBalance.Amount)

		balance = balance.Shift(-decoder.wm.Decimal())

		if balance.LessThan(minTransfer) || balance.LessThanOrEqual(decimal.Zero) {
			continue
		}
		//计算汇总数量 = 余额 - 保留余额
		sumAmount := balance.Sub(retainedBalance)

		//减去手续费
		sumAmount = sumAmount.Sub(fee)
		if sumAmount.LessThanOrEqual(decimal.Zero) {
			continue
		}

		decoder.wm.Log.Debugf("balance: %v", balance.String())
		decoder.wm.Log.Debugf("fees: %v", fee.String())
		decoder.wm.Log.Debugf("sumAmount: %v", sumAmount.String())

		//创建一笔交易单
		rawTx := &openwallet.RawTransaction{
			Coin:    sumRawTx.Coin,
			Account: sumRawTx.Account,
			To: map[string]string{
				sumRawTx.SummaryAddress: sumAmount.String(),
			},
			Required: 1,
			Fees:fee.String(),
		}

		createErr := decoder.createRawTransaction(
			wrapper,
			rawTx,
			addrBalance)
		if createErr != nil {
			return nil, createErr
		}

		//创建成功，添加到队列
		rawTxArray = append(rawTxArray, rawTx)

	}

	return rawTxArray, nil

}

//createRawTransaction
func (decoder *TransactionDecoder) createRawTransaction(
	wrapper openwallet.WalletDAI,
	rawTx *openwallet.RawTransaction,
	addrBalance *LTGAccount) error {

	var (
		accountTotalSent = decimal.Zero
		txFrom           = make([]string, 0)
		txTo             = make([]string, 0)
		keySignList      = make([]*openwallet.KeySignature, 0)
		amountStr        string
		destination      string
	)

	decimals := decoder.wm.Decimal()

	for k, v := range rawTx.To {
		destination = k
		amountStr = v
		break
	}

	amountTo, _ := decimal.NewFromString(amountStr)
	amountTo = amountTo.Shift(decimals)

	addr, err := wrapper.GetAddress(addrBalance.Address)
	if err != nil {
		return err
	}

	txFrom = []string{fmt.Sprintf("%s:%s", addr.Address, amountStr)}
	txTo = []string{fmt.Sprintf("%s:%s", destination, amountStr)}

	//nonce := decoder.wm.GetAddressNonce(wrapper, addrBalance)

	//decoder.wm.Log.Debugf("nonce: %d", nonce)
	//nonce = nonce + 1

	eop := transaction.EmptyOrdinaryPayment()

	pub, _ := hex.DecodeString(addr.PublicKey)

	enpub, _ := owcrypt.CURVE25519_convert_Ed_to_X(pub)

	feeToDecimal,_ := decimal.NewFromString(rawTx.Fees)

	feeTo := feeToDecimal.Shift(decimals)

	//目标LTG转换
	recipient, _ := rsencoding.Decode(destination[4:])

	block, err := decoder.wm.GetLatestBlock()
	if err != nil {
		return err
	}

	blockID, err := strconv.ParseUint(block.Hash, 10, 64)
	if err != nil {
		return err
	}

	memo := rawTx.GetExtParam().Get("memo").String()

	if memo != "" {
		eop.Appendix = &p2p.Appendix{
			Message: &p2p.Appendix_Message{
				Content: []byte(memo),
				IsText:  true,
			},
		}
	}

	eop.Header = &p2p.TransactionHeader{
		SenderPublicKey: enpub,
		Recipient:       recipient,
		Amount:          uint64(amountTo.IntPart()),
		Fee:             uint64(feeTo.IntPart()),
		Deadline:        360,
		Timestamp:       uint32(time.Now().Unix() - 1561852800),
		Version:         1,
		EcBlockHeight:   uint32(block.Height),
		EcBlockId:       uint64(blockID),
	}

	txStr := hex.EncodeToString(transaction.ToBytes(eop))

	tx := &RawTransactionV2{
		Sender:        hex.EncodeToString(enpub),
		Recipient:     eop.Header.Recipient,
		Symbol:        rawTx.Coin.Symbol,
		Amount:        eop.Header.Amount,
		Fee:           eop.Header.Fee,
		DeadLine:      eop.Header.Deadline,
		Timestamp:     eop.Header.Timestamp,
		Version:       eop.Header.Version,
		EcBlockHeight: eop.Header.EcBlockHeight,
		EcBlockId:     eop.Header.EcBlockId,
	}
	if memo != "" {
		tx.Appendix = &SendAppendix{
			Message:       memo,
			MessageIsText: true,
			Version:       1,
		}
	}

	txJson, _ := json.Marshal(tx)
	rawTx.RawHex = hex.EncodeToString(txJson)

	//rawTx.RawHex = txStr

	if rawTx.Signatures == nil {
		rawTx.Signatures = make(map[string][]*openwallet.KeySignature)
	}

	signature := openwallet.KeySignature{
		EccType: decoder.wm.CurveType(),
		Address: addr,
		Message: txStr,
	}
	keySignList = append(keySignList, &signature)

	accountTotalSent = decimal.Zero.Sub(accountTotalSent)
	rawTx.Signatures[rawTx.Account.AccountID] = keySignList
	rawTx.FeeRate = ""
	rawTx.IsBuilt = true
	rawTx.TxAmount = accountTotalSent.StringFixed(decimals)
	rawTx.TxFrom = txFrom
	rawTx.TxTo = txTo

	return nil
}

//CreateSummaryRawTransactionWithError 创建汇总交易，返回能原始交易单数组（包含带错误的原始交易单）
func (decoder *TransactionDecoder) CreateSummaryRawTransactionWithError(wrapper openwallet.WalletDAI, sumRawTx *openwallet.SummaryRawTransaction) ([]*openwallet.RawTransactionWithError, error) {
	raTxWithErr := make([]*openwallet.RawTransactionWithError, 0)
	rawTxs, err := decoder.CreateSummaryRawTransaction(wrapper, sumRawTx)
	if err != nil {
		return nil, err
	}
	for _, tx := range rawTxs {
		raTxWithErr = append(raTxWithErr, &openwallet.RawTransactionWithError{
			RawTx: tx,
			Error: nil,
		})
	}
	return raTxWithErr, nil
}
