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
	"fmt"
	"github.com/blocktree/openwallet/v2/log"
	"github.com/blocktree/openwallet/v2/openw"
	"github.com/blocktree/openwallet/v2/openwallet"
	"testing"
	"time"
)

func testGetAssetsAccountBalance(tm *openw.WalletManager, walletID, accountID string) {
	balance, err := tm.GetAssetsAccountBalance(testApp, walletID, accountID)
	if err != nil {
		log.Error("GetAssetsAccountBalance failed, unexpected error:", err)
		return
	}
	log.Info("balance:", balance)
}

func testGetAssetsAccountTokenBalance(tm *openw.WalletManager, walletID, accountID string, contract openwallet.SmartContract) {
	balance, err := tm.GetAssetsAccountTokenBalance(testApp, walletID, accountID, contract)
	if err != nil {
		log.Error("GetAssetsAccountTokenBalance failed, unexpected error:", err)
		return
	}
	log.Info("token balance:", balance.Balance)
}

func testCreateTransactionStep(tm *openw.WalletManager, walletID, accountID, to, amount, feeRate,memo string, contract *openwallet.SmartContract) (*openwallet.RawTransaction, error) {

	//err := tm.RefreshAssetsAccountBalance(testApp, accountID)
	//if err != nil {
	//	log.Error("RefreshAssetsAccountBalance failed, unexpected error:", err)
	//	return nil, err
	//}

	rawTx, err := tm.CreateTransaction(testApp, walletID, accountID, amount, to, feeRate, memo, contract)

	if err != nil {
		log.Error("CreateTransaction failed, unexpected error:", err)
		return nil, err
	}

	return rawTx, nil
}

func testCreateSummaryTransactionStep(
	tm *openw.WalletManager,
	walletID, accountID, summaryAddress, minTransfer, retainedBalance, feeRate string,
	start, limit int,
	contract *openwallet.SmartContract) ([]*openwallet.RawTransaction, error) {

	rawTxArray, err := tm.CreateSummaryTransaction(testApp, walletID, accountID, summaryAddress, minTransfer,
		retainedBalance, feeRate, start, limit, contract)

	if err != nil {
		log.Error("CreateSummaryTransaction failed, unexpected error:", err)
		return nil, err
	}

	return rawTxArray, nil
}

func testSignTransactionStep(tm *openw.WalletManager, rawTx *openwallet.RawTransaction) (*openwallet.RawTransaction, error) {

	_, err := tm.SignTransaction(testApp, rawTx.Account.WalletID, rawTx.Account.AccountID, "12345678", rawTx)
	if err != nil {
		log.Error("SignTransaction failed, unexpected error:", err)
		return nil, err
	}

	log.Infof("rawTx: %+v", rawTx)
	return rawTx, nil
}

func testVerifyTransactionStep(tm *openw.WalletManager, rawTx *openwallet.RawTransaction) (*openwallet.RawTransaction, error) {

	//log.Info("rawTx.Signatures:", rawTx.Signatures)

	_, err := tm.VerifyTransaction(testApp, rawTx.Account.WalletID, rawTx.Account.AccountID, rawTx)
	if err != nil {
		log.Error("VerifyTransaction failed, unexpected error:", err)
		return nil, err
	}

	log.Infof("rawTx: %+v", rawTx)
	return rawTx, nil
}

func testSubmitTransactionStep(tm *openw.WalletManager, rawTx *openwallet.RawTransaction) (*openwallet.RawTransaction, error) {

	tx, err := tm.SubmitTransaction(testApp, rawTx.Account.WalletID, rawTx.Account.AccountID, rawTx)
	if err != nil {
		log.Error("SubmitTransaction failed, unexpected error:", err)
		return nil, err
	}

	log.Std.Info("tx: %+v", tx)
	log.Info("wxID:", tx.WxID)
	log.Info("txID:", rawTx.TxID)

	return rawTx, nil
}

func TestTransfer(t *testing.T) {

	//LTG-45DW-CZQU-T24P-E8NUU
	//LTG-49PH-KJ4L-GTS5-64ZQU
	//LTG-4VSG-H2MR-Y2QK-47JWG
	//LTG-6C62-K7KG-FDLG-FL7F5
	//LTG-73XT-BS88-QN5R-9P2KQ
	//LTG-8HCX-64UF-BA37-9JKV5
	//LTG-9QDM-VWFN-ZFRV-B5HF6
	//LTG-AY35-82UW-DFJ3-BT42A
	//LTG-GSFA-V7XY-3VR9-AE9ND
	//LTG-RFW7-USTF-66C7-E3A6K
	//LTG-RPYP-Y4R2-2W5A-GLTWU

	addrs := []string{
		//"LTG-8JKU-8ZGQ-3LRS-9NPVH",
		//"LTG-73XT-BS88-QN5R-9P2KQ",
		"LTG-45DW-CZQU-T24P-E8NUU",
		"LTG-49PH-KJ4L-GTS5-64ZQU",
		"LTG-RFW7-USTF-66C7-E3A6K",
		"LTG-6C62-K7KG-FDLG-FL7F5",
		"LTG-AY35-82UW-DFJ3-BT42A",
		"LTG-GSFA-V7XY-3VR9-AE9ND",
		"LTG-8HCX-64UF-BA37-9JKV5",
		//"LTG-8JKU-8ZGQ-3LRS-9NPVH",
		//"LTG-K6UC-JHM2-V627-H4W34",
		//"LTG-TMLD-S5TB-78Y2-2KLXK",
		//"LTG-VEY6-VWBX-QGSU-FSAJ9",
		//"LTG-NXZF-TQ4U-H7RV-G24BE",
		//"LTG-UTHN-X6AR-ZKTR-DWVLM",
		//"LTG-E5E3-D8W7-3BE2-DEWVD",
		//"LTG-HUJ6-7K7G-SP4T-8J6CP",
		//"LTG-86X2-W683-S2FN-2KKUT",
	}

	tm := testInitWalletManager()
	//walletID := "WAVP3DMctqqLVe8owHktCekfGon24Xg34r"
	//accountID := "GRTCmSZoczPPkxGH1NJKfwYJerJHGtGP2TabeTLaBR6m"
	//walletID := "W33DGqQ3prBXh2qXqwJ2YZiATpTxYZPsEs"
	//accountID := "sVEMydRNKc8hJiug3jiiC1PmY9mM2vdMZsaAFMHLomC"

	walletID := "WJCbRxzKfukURWonxTPd7rxvUbVBEgRBMW"
	accountID := "DdFbucjCkNopZTTvLKWgRDscJnaMiV49eBwUdtn7mQ4T"
	testGetAssetsAccountBalance(tm, walletID, accountID)
	for _, to := range addrs {

		rawTx, err := testCreateTransactionStep(tm, walletID, accountID, to, "0.1", "", "998877771",nil)
		if err != nil {
			return
		}

		_, err = testSignTransactionStep(tm, rawTx)
		if err != nil {
			return
		}

		_, err = testVerifyTransactionStep(tm, rawTx)
		if err != nil {
			return
		}

		_, err = testSubmitTransactionStep(tm, rawTx)
		if err != nil {
			return
		}
		fmt.Print("fees:")
		fmt.Println(rawTx.Fees)

		time.Sleep(5 * time.Second)
	}
}

func TestSummary(t *testing.T) {
	tm := testInitWalletManager()
	//walletID := "WAVP3DMctqqLVe8owHktCekfGon24Xg34r"
	//accountID := "GRTCmSZoczPPkxGH1NJKfwYJerJHGtGP2TabeTLaBR6m"

	walletID := "WLobhbspGbPD6GMeyVzvKwWHKnff23T81p"
	accountID := "6HduTD8nqcWKWuoFbCJoWX7x9JB683XPUp1dNRFSYD1M"
	summaryAddress := "LTG-NKUN-SFN9-4D4U-7X8UE"

	testGetAssetsAccountBalance(tm, walletID, accountID)

	rawTxArray, err := testCreateSummaryTransactionStep(tm, walletID, accountID,
		summaryAddress, "", "", "",
		0, 100, nil)
	if err != nil {
		log.Errorf("CreateSummaryTransaction failed, unexpected error: %v", err)
		return
	}

	//执行汇总交易
	for _, rawTx := range rawTxArray {
		_, err = testSignTransactionStep(tm, rawTx)
		if err != nil {
			return
		}

		_, err = testVerifyTransactionStep(tm, rawTx)
		if err != nil {
			return
		}

		_, err = testSubmitTransactionStep(tm, rawTx)
		if err != nil {
			return
		}
	}

}
