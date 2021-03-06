/*
 * Copyright 2018 The openwallet Authors
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

package openwtester

import (
	"github.com/astaxie/beego/config"
	"github.com/blocktree/openwallet/v2/common/file"
	"github.com/blocktree/openwallet/v2/log"
	"github.com/blocktree/openwallet/v2/openw"
	"github.com/blocktree/openwallet/v2/openwallet"
	"path/filepath"
	"testing"
)

////////////////////////// 测试单个扫描器 //////////////////////////

type subscriberSingle struct {
}

//BlockScanNotify 新区块扫描完成通知
func (sub *subscriberSingle) BlockScanNotify(header *openwallet.BlockHeader) error {
	log.Std.Notice("header: %+v", header)
	return nil
}

//BlockTxExtractDataNotify 区块提取结果通知
func (sub *subscriberSingle) BlockExtractDataNotify(sourceKey string, data *openwallet.TxExtractData) error {
	log.Notice("account:", sourceKey)
	log.Notice("fee:", data.Transaction.Fees)

	if data.Transaction.GetExtParam().Get("memo").String() != ""{
		log.Std.Notice("memo is [%s]",data.Transaction.GetExtParam().Get("memo").String())
	}else{
		log.Std.Notice("memo is nil")
	}

	for i, input := range data.TxInputs {
		log.Std.Notice("data.TxInputs[%d]: %+v", i, input)

		log.Notice("intput:to:",input.Amount)
	}

	for i, output := range data.TxOutputs {
		log.Std.Notice("data.TxOutputs[%d]: %+v", i, output)
		log.Notice("output:to:",output.Amount)

	}

	log.Std.Notice("data.Transaction: %+v", data.Transaction)

	return nil
}


func (sub *subscriberSingle) BlockExtractSmartContractDataNotify(sourceKey string, data *openwallet.SmartContractReceipt) error {

	log.Notice("sourceKey:", sourceKey)
	log.Std.Notice("data.ContractTransaction: %+v", data)

	for i, event := range data.Events {
		log.Std.Notice("data.Events[%d]: %+v", i, event)
	}

	return nil
}

func TestSubscribeAddress(t *testing.T) {

	var (
		endRunning = make(chan bool, 1)
		symbol     = "LTG"
		addrs      = map[string]string{
			"LTG-NKUN-SFN9-4D4U-7X8UE": "sender",
			"LTG-45DW-CZQU-T24P-E8NUU": "receiver",
			"LTG-49PH-KJ4L-GTS5-64ZQU": "receiver",
			"LTG-RFW7-USTF-66C7-E3A6K": "receiver",
			"LTG-6C62-K7KG-FDLG-FL7F5": "receiver",
			"LTG-AY35-82UW-DFJ3-BT42A": "receiver",
			"LTG-GSFA-V7XY-3VR9-AE9ND": "receiver",
			"LTG-8HCX-64UF-BA37-9JKV5": "receiver",
			"LTG-VNKD-C65H-QAZF-B5CAZ":"receiver",
		}
	)

	//GetSourceKeyByAddress 获取地址对应的数据源标识
	scanTargetFunc := func(target openwallet.ScanTarget) (string, bool) {
		key, ok := addrs[target.Address]
		if !ok {
			return "", false
		}
		return key, true
	}

	assetsMgr, err := openw.GetAssetsAdapter(symbol)
	if err != nil {
		log.Error(symbol, "is not support")
		return
	}

	//读取配置
	absFile := filepath.Join(configFilePath, symbol+".ini")

	c, err := config.NewConfig("ini", absFile)
	if err != nil {
		return
	}
	assetsMgr.LoadAssetsConfig(c)

	assetsLogger := assetsMgr.GetAssetsLogger()
	if assetsLogger != nil {
		assetsLogger.SetLogFuncCall(true)
	}

	//log.Debug("already got scanner:", assetsMgr)
	scanner := assetsMgr.GetBlockScanner()

	if scanner.SupportBlockchainDAI() {
		file.MkdirAll(dbFilePath)
		dai, err := openwallet.NewBlockchainLocal(filepath.Join(dbFilePath, dbFileName), false)
		if err != nil {
			log.Error("NewBlockchainLocal err: %v", err)
			return
		}

		scanner.SetBlockchainDAI(dai)
	}

	scanner.SetRescanBlockHeight(39179)

	if scanner == nil {
		log.Error(symbol, "is not support block scan")
		return
	}

	scanner.SetBlockScanTargetFunc(scanTargetFunc)

	sub := subscriberSingle{}
	scanner.AddObserver(&sub)

	scanner.Run()

	<-endRunning
}
