/*
 * Copyright 2018 The OpenWallet Authors
 * This file is part of the OpenWallet library.
 *
 * The OpenWallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The OpenWallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package livingthegrace

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/assetsadapterstore/livingthegrace-adapter/livingthegrace_addrdec"
	"github.com/blocktree/go-owcrypt"
	"github.com/blocktree/openwallet/v2/common"
	"github.com/blocktree/openwallet/v2/log"
	"github.com/blocktree/openwallet/v2/openwallet"
	"github.com/imroc/req"
	"github.com/tidwall/gjson"
)

type WalletManager struct {
	openwallet.AssetsAdapterBase

	client       *Client                       // 节点客户端
	Config       *WalletConfig                 // 节点配置
	Decoder      openwallet.AddressDecoder     //地址编码器
	DecoderV2    openwallet.AddressDecoderV2   //地址编码器V2
	TxDecoder    openwallet.TransactionDecoder //交易单编码器
	Log          *log.OWLogger                 //日志工具
	Blockscanner openwallet.BlockScanner       //区块扫描器
}

func NewWalletManager() *WalletManager {
	wm := WalletManager{}
	wm.Config = NewConfig(Symbol)
	wm.Blockscanner = NewBlockScanner(&wm)
	wm.Decoder = NewAddressDecoder()
	wm.DecoderV2 = &livingthegrace_addrdec.AddressDecoderV2{}
	wm.TxDecoder = NewTransactionDecoder(&wm)
	wm.Log = log.NewOWLogger(wm.Symbol())
	return &wm
}

func (wm *WalletManager) GetWalletDetails(address string) (*LTGAccount, error) {
	path := fmt.Sprintf("/burst?requestType=getBalance&account=%s", address)

	result, err := wm.client.call("GET", path, nil)
	if err != nil {
		return nil, err
	}

	return NewLTGAccount(address, result), nil
}

func (wm *WalletManager) NewWallet(symbol string) (*gjson.Result, error) {

	path := fmt.Sprintf("coin/new")

	pararm := req.Param{
		"symbol": symbol,
	}

	result, err := wm.client.call("POST", path, pararm)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (wm *WalletManager) GetLatestBlock() (*Block, error) {

	path := fmt.Sprintf("/burst?requestType=getBlock")
	result, err := wm.client.call("GET", path, nil)
	if err != nil {
		return nil, err
	}
	return NewBlock(result), nil
}


func (wm *WalletManager) GetSupportFee() (*SupportFee, error) {

	path := fmt.Sprintf("/burst?requestType=suggestFee")
	result, err := wm.client.call("GET", path, nil)
	if err != nil {
		return nil, err
	}
	return NewSupportFee(result), nil
}

func (wm *WalletManager) GetBlock(num uint64) (*Block, error) {

	path := fmt.Sprintf("/burst?requestType=getBlock&height=%d", num)
	result, err := wm.client.call("GET", path, nil)
	if err != nil {
		return nil, err
	}
	return NewBlock(result), nil
}

func (wm *WalletManager) GetTransaction(txid string) (*Transaction, error) {

	path := fmt.Sprintf("burst?requestType=getTransaction&transaction=%s", txid)
	result, err := wm.client.call("GET", path, nil)
	if err != nil {
		return nil, err
	}

	if result.Get("errorCode").Exists() {
		return nil, errors.New(result.Get("errorDescription").String())
	}

	return NewTransaction(result), nil
}

func (wm *WalletManager) Sendraw(rawTx *RawTransaction) (string, error) {

	path := fmt.Sprintf("coin/sendraw")

	pararm := req.Param{
		"sender":    rawTx.Sender,
		"recipient": rawTx.Recipient,
		"symbol":    rawTx.Symbol,
		"amount":    rawTx.Amount,
		"nonce":     rawTx.Nonce,
		"signature": rawTx.Signature,
	}

	result, err := wm.client.call("POST", path, pararm)
	if err != nil {
		return "", err
	}

	return result.Get("txn").String(), nil
}

func (wm *WalletManager) SendOffTx(tx string) (string, error) {

	path := fmt.Sprintf("burst?requestType=broadcastTransaction")

	pararm := req.Param{
		"transactionJSON": tx,
	}
	result, err := wm.client.call("POST", path, pararm)
	if err != nil {
		return "", err
	}

	if result.Get("errorCode").Exists() {
		return "", errors.New(result.Get("errorDescription").String())
	}

	fmt.Println(result.Raw)

	return result.Get("transaction").String(), nil
}

func (wm *WalletManager) SignRawTxOnline(rawTx *RawTransaction, privateKey []byte) error {

	path := fmt.Sprintf("signature/generate")

	pararm := req.Param{
		"sender":     rawTx.Sender,
		"recipient":  rawTx.Recipient,
		"symbol":     rawTx.Symbol,
		"amount":     rawTx.Amount,
		"nonce":      rawTx.Nonce,
		"privatekey": hex.EncodeToString(privateKey),
	}

	result, err := wm.client.call("POST", path, pararm)
	if err != nil {
		return err
	}
	rawTx.Signature = result.Get("signature").String()
	return nil
}

func (wm *WalletManager) SignRawTxOffline(rawTx *RawTransaction, privateKey []byte) error {

	messageHash := rawTx.Hash()
	signature, _, ret := owcrypt.Signature(privateKey, nil, messageHash, wm.CurveType())
	if ret != owcrypt.SUCCESS {
		return fmt.Errorf("sign raw tx failed")
	}
	rawTx.FillSig(signature)
	publicKey, _ := hex.DecodeString(rawTx.Sender)
	pub := owcrypt.PointDecompress(publicKey, wm.CurveType())
	verRet := owcrypt.Verify(pub[1:], nil, messageHash, signature, wm.CurveType())
	if verRet != owcrypt.SUCCESS {
		log.Errorf("transaction verify failed")
	} else {
		log.Infof("transaction verify success")
	}
	return nil
}

// GetAddressNonce
func (wm *WalletManager) GetAddressNonce(wrapper openwallet.WalletDAI, account *LTGAccount) uint64 {
	var (
		//key           = wm.Symbol() + "-nonce"
		nonce         uint64
		nonce_db      interface{}
		nonce_onchain uint64
	)

	//获取db记录的nonce并确认nonce值
	//nonce_db, _ = wrapper.GetAddressExtParam(account.Publickey, key)

	//判断nonce_db是否为空,为空则说明当前nonce是0
	if nonce_db == nil {
		nonce = 0
	} else {
		nonce = common.NewString(nonce_db).UInt64()
	}

	//nonce_onchain = account.Nonce

	//如果本地nonce_db > 链上nonce,采用本地nonce,否则采用链上nonce
	if nonce > nonce_onchain {
		//wm.Log.Debugf("%s nonce_db=%v > nonce_chain=%v,Use nonce_db...", address, nonce_db, nonce_onchain)
	} else {
		nonce = nonce_onchain
		//wm.Log.Debugf("%s nonce_db=%v <= nonce_chain=%v,Use nonce_chain...", address, nonce_db, nonce_onchain)
	}

	return nonce
}

// UpdateAddressNonce
func (wm *WalletManager) UpdateAddressNonce(wrapper openwallet.WalletDAI, address string, nonce uint64) {
	key := wm.Symbol() + "-nonce"
	err := wrapper.SetAddressExtParam(address, key, nonce)
	if err != nil {
		wm.Log.Errorf("WalletDAI SetAddressExtParam failed, err: %v", err)
	}
}
