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
	"github.com/blocktree/go-owcrypt"
	"github.com/shopspring/decimal"
)

const (
	CurveType = owcrypt.ECC_CURVE_CURVE25519_SHA256
	Symbol    = "LTG"
)

type WalletConfig struct {

	//币种
	Symbol string
	//配置文件路径
	configFilePath string
	//配置文件名
	configFileName string
	//钱包服务API
	ServerAPI string
	//曲线类型
	CurveType uint32
	//数据目录
	DataDir string
	//Fix Required Fee
	FixFees decimal.Decimal
}

func NewConfig(symbol string) *WalletConfig {

	c := WalletConfig{}

	//币种
	c.Symbol = symbol
	c.CurveType = CurveType
	//钱包服务API
	c.ServerAPI = ""

	c.FixFees,_ = decimal.NewFromString("0.073500")

	return &c
}
