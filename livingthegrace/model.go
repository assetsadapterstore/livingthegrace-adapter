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
	"fmt"
	"github.com/blocktree/go-owcrypt"
	"github.com/blocktree/openwallet/v2/openwallet"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
)

const TimeLayout = `2006-01-02T15:04:05Z07:00`
const INITTIME = 1561852800

type LTGAccount struct {
	Address               string `json:"address"`
	Symbol                string `json:"symbol"`
	Amount                string `json:"balanceNQT"`
	UnconfirmedBalance    string `json:"unconfirmedBalanceNQT"`
	EffectiveBalance      string `json:"effectiveBalanceNXT"`
	GuaranteedBalanceNQT  string `json:"guaranteedBalanceNQT"`
	ForgedBalanceNQT      string `json:"forgedBalanceNQT"`
	RequestProcessingTime uint64 `json:"requestProcessingTime"`
	Nonce                 uint64
	Publickey             string
}

func NewLTGAccount(address string, result *gjson.Result) *LTGAccount {
	obj := LTGAccount{}
	obj.Address = address
	obj.Amount = result.Get("balanceNQT").String()
	obj.UnconfirmedBalance = result.Get("unconfirmedBalanceNQT").String()
	obj.EffectiveBalance = result.Get("effectiveBalanceNXT").String()
	obj.GuaranteedBalanceNQT = result.Get("guaranteedBalanceNQT").String()
	obj.ForgedBalanceNQT = result.Get("forgedBalanceNQT").String()
	obj.RequestProcessingTime = result.Get("requestProcessingTime").Uint()
	return &obj
}

type Block struct {
	Height    uint64   `json:"height"`
	Hash      string   `json:"block"`
	LastHash  string   `json:"previousBlock"`
	Txns      []string `json:"transactions"`
	Timestamp uint64   `json:"timestamp"`
}

func NewBlock(result *gjson.Result) *Block {
	obj := Block{}
	obj.Height = result.Get("height").Uint()
	obj.Hash = result.Get("block").String()
	obj.LastHash = result.Get("previousBlock").String()
	obj.Timestamp = result.Get("timestamp").Uint() + 1561852800
	obj.Txns = make([]string, 0)
	if txns := result.Get("transactions"); txns.IsArray() {
		for _, tx := range txns.Array() {
			obj.Txns = append(obj.Txns, tx.String())
		}
	}
	return &obj
}

//BlockHeader 区块链头
func (b *Block) BlockHeader(symbol string) *openwallet.BlockHeader {

	obj := openwallet.BlockHeader{}
	//解析json
	obj.Hash = b.Hash
	obj.Previousblockhash = b.LastHash
	obj.Height = b.Height
	obj.Time = b.Timestamp
	obj.Symbol = symbol

	return &obj
}

type Transaction struct {
	/*
		{
			"type": 0,
			"subtype": 0,
			"timestamp": 22922438,
			"deadline": 1440,
			"senderPublicKey": "542b31dfe9eeda053781ecc676f6c2a3e06cebc4c3c57446b82f60eb1b874123",
			"recipient": "11502548938569667115",
			"recipientRS": "LTG-VNKD-C65H-QAZF-B5CAZ",
			"amountNQT": "114002205000",
			"feeNQT": "100000000",
			"signature": "7741dbd4ce676688a7e7341a7c3bbe6f9ea5a6ae0d5c2d49ed5913cf65b7990e238ecf60d5314811ebfd519292a2958c5362fcc44b051d0bd405b07ac6671d2d",
			"signatureHash": "3241711592a27d58f67328092592093c2cb62626629676b85ad97cd57c4be489",
			"fullHash": "6378d06a707f3478d90398039bc4bc1843b1c8680dac7d50a8d93480ec2c2ca3",
			"transaction": "8661688104145418339",
			"sender": "10222100220764572315",
			"senderRS": "LTG-ULNV-ZRLK-9LZQ-AK3RV",
			"height": 33089,
			"version": 1,
			"ecBlockId": "9377985466004668487",
			"ecBlockHeight": 33076,
			"block": "3062415563766137057",
			"confirmations": 1,
			"blockTimestamp": 22922563,
			"requestProcessingTime": 2
		}
	*/
	Hash         string
	From         string
	To           string
	Amount       string
	Fee          string
	Symbol       string
	BlockHash    string
	BlockHeight  uint64
	Status       string
	TxType       string
	SubType      string
	Memo         string
	Timestamp    uint64
	SendAppendix *SendAppendix
}

func NewTransaction(result *gjson.Result) *Transaction {
	obj := Transaction{}
	obj.Hash = result.Get("transaction").String()
	obj.From = result.Get("senderRS").String()
	obj.To = result.Get("recipientRS").String()
	obj.Amount = result.Get("amountNQT").String()
	amountDec, _ := decimal.NewFromString(obj.Amount)
	amountDec = amountDec.Shift(-8)
	obj.Amount = amountDec.String()
	obj.Symbol = "LTG"
	obj.BlockHash = result.Get("block").String()
	obj.BlockHeight = result.Get("height").Uint()

	obj.Fee = result.Get("feeNQT").String()
	feeDec, _ := decimal.NewFromString(obj.Fee)
	feeDec = feeDec.Shift(-8)
	obj.Fee = feeDec.String()

	obj.Status = "1"
	obj.TxType = result.Get("type").String()
	obj.SubType = result.Get("subtype").String()
	obj.Timestamp = result.Get("timestamp").Uint() + 1561852800
	if result.Get("attachment").Exists() {
		obj.SendAppendix = &SendAppendix{
			Version:       int64(result.Get("attachment").Get("version\\.Message").Uint()),
			Message:       result.Get("attachment").Get("message").String(),
			MessageIsText: result.Get("attachment").Get("messageIsText").Bool(),
		}
	}
	if obj.SendAppendix != nil {
		if obj.SendAppendix.Version == 1 && obj.SendAppendix.MessageIsText {
			obj.Memo = obj.SendAppendix.Message
		}
	}
	return &obj
}

type RawTransaction struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Symbol    string `json:"symbol"`
	Amount    string `json:"amount"`
	Nonce     uint64 `json:"nonce"`
	Signature string `json:"signature"`
}

type RawTransactionV2 struct {
	Sender        string        `json:"sender"`
	Recipient     uint64        `json:"recipient"`
	Symbol        string        `json:"symbol"`
	DeadLine      uint32        `json:"deadLine"`
	Fee           uint64        `json:"fee"`
	Amount        uint64        `json:"amount"`
	Signature     string        `json:"signature"`
	Timestamp     uint32        `json:"timestamp"`
	EcBlockHeight uint32        `json:"ecBlockHeight,omitempty"`
	EcBlockId     uint64        `json:"ecBlockId,omitempty"`
	Version       uint32        `json:"version,omitempty"`
	Appendix      *SendAppendix `json:"attachment"`
}

type RawTransactionSend struct {
	EcBlockHeight   uint32        `json:"ecBlockHeight,omitempty"`
	EcBlockId       uint64        `json:"ecBlockId,omitempty"`
	Version         uint32        `json:"version,omitempty"`
	Type            uint64        `json:"type"`
	Deadline        uint32        `json:"deadline"`
	AmountNQT       uint64        `json:"amountNQT"`
	FeeNQT          uint64        `json:"feeNQT"`
	Signature       string        `json:"signature"`
	SenderPublicKey string        `json:"senderPublicKey"`
	Recipient       uint64        `json:"recipient"`
	Timestamp       uint32        `json:"timestamp"`
	Appendix        *SendAppendix `json:"attachment,omitempty"`
}

type SendAppendix struct {
	Version       int64  `json:"version.Message"`
	Message       string `json:"message"`
	MessageIsText bool   `json:"messageIsText"`
}

func (rawTx *RawTransaction) Hash() []byte {
	message := fmt.Sprintf("%s%s%s%s%d", rawTx.Sender, rawTx.Recipient, rawTx.Symbol, rawTx.Amount, rawTx.Nonce)
	messageHash := owcrypt.Hash([]byte(message), 0, owcrypt.HASH_ALG_SHA256)
	return messageHash
}

func (rawTx *RawTransaction) FillSig(signature []byte) error {
	if len(signature) != 64 {
		return fmt.Errorf("signature length is not equal 64 bytes")
	}
	//DER-encoded, 30440220+前32字节+0220+后32字节
	lBytes := signature[:32]
	rBytes := signature[32:]
	der := append([]byte{0x30, 0x44, 0x02, 0x20}, lBytes...)
	der = append(der, 0x02, 0x20)
	der = append(der, rBytes...)
	rawTx.Signature = hex.EncodeToString(der)
	return nil
}
