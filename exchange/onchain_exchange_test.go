/*
 * Copyright (C) 2020 The poly network Authors
 * This file is part of The poly network library.
 *
 * The  poly network  is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The  poly network  is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 * You should have received a copy of the GNU Lesser General Public License
 * along with The poly network .  If not, see <http://www.gnu.org/licenses/>.
 */

package exchange

import (
	"encoding/hex"
	"github.com/ontio/ontology/common"
	"github.com/skyinglyh1/uniswap_v1_test/log"
	"math/big"
	"testing"
)

func Test_AddLiquidity(t *testing.T) {
	for i := 0; i < len(testEnv.OnChainEState[0].Providers); i++ {
		providerAddr := testEnv.OnChainEState[0].Providers[i].Address
		log.Infof("account: %s, ontDBalance: %+v, tokenBalance: %+v, exchange 1 shareBalance: %+v\n", providerAddr.ToBase58(), testEnv.OntdBalance[providerAddr], testEnv.OnChainTState[0].Balances[providerAddr], testEnv.OnChainEState[0].ShareBalance[providerAddr])
	}
	providerAddr := testEnv.OnChainEState[0].Providers[0].Address
	log.Infof("account: %s, ontDBalance: %+v, tokenBalance: %+v, shareBalance: %+v\n", providerAddr.ToBase58(), testEnv.OntdBalance[providerAddr], testEnv.OnChainTState[0].Balances[providerAddr], testEnv.OnChainEState[0].ShareBalance[providerAddr])

	if err := testEnv.addLiquid(0, big.NewInt(1), big.NewInt(0).Add(big.NewInt(100000), big.NewInt(10000000)), big.NewInt(200000)); err != nil {
		log.Errorf("address: %s, addLiquid() error: %+v", err)
	}
	if err := testEnv.addLiquid(1, big.NewInt(1), big.NewInt(0).Add(big.NewInt(100000), big.NewInt(10000000)), big.NewInt(200000)); err != nil {
		log.Errorf("address: %s, addLiquid() error: %+v", err)
	}
}

func Test_RemoveLiquidity(t *testing.T) {
	providerAddr := testEnv.OnChainEState[0].Providers[0].Address
	log.Infof("account: %s, ontdBalance: %+v, tokenBalance: %+v, shareBalance: %+v\n", providerAddr.ToBase58(), testEnv.OntdBalance[providerAddr], testEnv.OnChainEState[0].ShareBalance[providerAddr], testEnv.OnChainEState[0].ShareBalance[providerAddr])

	if err := testEnv.removeLiquid(0, big.NewInt(1000), big.NewInt(1), testEnv.OnChainEState[0].Providers[0]); err != nil {
		log.Errorf("address: %s, removeLiquid() error: %+v", err)
	}
	if err := testEnv.removeLiquid(1, big.NewInt(1000), big.NewInt(1), testEnv.OnChainEState[0].Providers[0]); err != nil {
		log.Errorf("address: %s, removeLiquid() error: %+v", err)
	}
}

func Test_ontToTokenInput(t *testing.T) {
	usrAddr := testEnv.Users[0].Address
	log.Infof("account: %s, ontdBalance: %+v, tokenBalance: %+v, shareBalance: %+v\n", usrAddr.ToBase58(), testEnv.OntdBalance[usrAddr], testEnv.OnChainEState[0].ShareBalance[usrAddr], testEnv.OnChainEState[0].ShareBalance[usrAddr])
	ontdSold := big.NewInt(10)
	minTokens := big.NewInt(1)

	testEnv.OnChainEState[0].offOntToTokenInput(big.NewInt(5), minTokens)
	if err := testEnv.ontToTokenInput(ontdSold, minTokens, testEnv.OnChainEState[0].Providers[0], testEnv.OnChainEState[0].Providers[0].Address); err != nil {
		log.Errorf("address: %s, ongToTokenSwapInput() error: %+v", err)
	}
	if err := testEnv.ontToTokenInput(ontdSold, minTokens, testEnv.OnChainEState[0].Providers[0], testEnv.Users[1].Address); err != nil {
		log.Errorf("address: %s, ongToTokenTransferInput() error: %+v", err)
	}
}

func Test_ontToTokenOutput(t *testing.T) {
	usrAddr := testEnv.Users[0].Address
	log.Infof("account: %s, ontdBalance: %+v, tokenBalance: %+v, shareBalance: %+v\n", usrAddr.ToBase58(), testEnv.OntdBalance[usrAddr], testEnv.OnChainEState[0].ShareBalance[usrAddr], testEnv.OnChainEState[0].ShareBalance[usrAddr])

	tokenBought := big.NewInt(10)
	maxOntd := big.NewInt(100)

	testEnv.OnChainEState[0].offOntToTokenOutput(big.NewInt(5), maxOntd)
	if err := testEnv.ontToTokenOutput(tokenBought, maxOntd, testEnv.Users[0], testEnv.Users[0].Address); err != nil {
		log.Errorf("address: %s, ongToTokenSwapOutput() error: %+v", err)
	}
	if err := testEnv.ontToTokenOutput(tokenBought, maxOntd, testEnv.Users[0], testEnv.Users[1].Address); err != nil {
		log.Errorf("address: %s, ongToTokeTransferpOutput() error: %+v", err)
	}
}

func Test_tokenToOntInput(t *testing.T) {
	providerAddr := testEnv.OnChainEState[0].Providers[0].Address
	log.Infof("account: %s, ontdBalance: %+v, tokenBalance: %+v, shareBalance: %+v\n", providerAddr.ToBase58(), testEnv.OntdBalance[providerAddr], testEnv.OnChainEState[0].ShareBalance[providerAddr], testEnv.OnChainEState[0].ShareBalance[providerAddr])

	tokenSold := big.NewInt(5)
	minOng := big.NewInt(1)

	testEnv.OnChainEState[0].offTokenToOntInput(big.NewInt(5), minOng)
	if err := testEnv.tokenToOntInput(tokenSold, minOng, testEnv.Users[0], testEnv.Users[0].Address); err != nil {
		log.Errorf("address1: %s, tokenToOngSwapInput() error: %+v", err)
	}
	if err := testEnv.tokenToOntInput(tokenSold, minOng, testEnv.Users[0], testEnv.Users[1].Address); err != nil {
		log.Errorf("address2: %s, tokenToOngTransferInput() error: %+v", err)
	}
}

func Test_tokenToOntOutput(t *testing.T) {
	usrAddr := testEnv.Users[0].Address
	log.Infof("account: %s, ontdBalance: %+v, tokenBalance: %+v, shareBalance: %+v\n", usrAddr.ToBase58(), testEnv.OntdBalance[usrAddr], testEnv.OnChainEState[0].ShareBalance[usrAddr], testEnv.OnChainEState[0].ShareBalance[usrAddr])

	var ongBought uint64 = 5
	maxTokens := big.NewInt(100)

	testEnv.OnChainEState[0].offTokenToOntOutput(big.NewInt(0).SetUint64(ongBought), maxTokens)
	if err := testEnv.tokenToOntOutput(ongBought, maxTokens, testEnv.Users[0], testEnv.Users[0].Address); err != nil {
		log.Errorf("address: %s, tokenToOngSwapInput() error: %+v", err)
	}
	if err := testEnv.tokenToOntOutput(ongBought, maxTokens, testEnv.Users[0], testEnv.Users[1].Address); err != nil {
		log.Errorf("address: %s, tokenToOngTransferInput() error: %+v", err)
	}
}

func Test_tokenToTokenInput(t *testing.T) {
	usrAddr := testEnv.Users[0].Address
	log.Infof("account: %s, ontdBalance: %+v, tokenBalance: %+v, shareBalance: %+v\n", usrAddr.ToBase58(), testEnv.OntdBalance[usrAddr], testEnv.OnChainEState[0].ShareBalance[usrAddr], testEnv.OnChainEState[0].ShareBalance[usrAddr])

	tokenSold := big.NewInt(50)

	ontdBought, tokenBought := testEnv.offTokenToTokenInput(tokenSold)
	minOntdBought, minTokenBought := ontdBought.Sub(ontdBought, big.NewInt(10)), tokenBought.Sub(tokenBought, big.NewInt(10))

	token1Hash, _ := common.AddressFromHexString(hex.EncodeToString(testEnv.OnChainTState[1].TokenAddr[:]))
	if err := testEnv.tokenToTokenInput(0, tokenSold, minTokenBought, minOntdBought, testEnv.OnChainEState[0].Providers[0], testEnv.Users[0].Address, token1Hash); err != nil {
		log.Errorf("address: %s, ongToTokenSwapInput() error: %+v", err)
	}
	ontdBought, tokenBought = testEnv.offTokenToTokenInput(tokenSold)
	minOntdBought, minTokenBought = ontdBought, tokenBought
	if err := testEnv.tokenToTokenInput(0, tokenSold, minTokenBought, minOntdBought, testEnv.OnChainEState[0].Providers[0], testEnv.Users[1].Address, token1Hash); err != nil {
		log.Errorf("address: %s, ongToTokenSwapInput() error: %+v", err)
	}
}

func Test_tokenToTokenOutput(t *testing.T) {
	userIndex := 4
	usrAddr := testEnv.Users[userIndex].Address
	log.Infof("account: %s, ontdBalance: %+v, tokenBalance: %+v, shareBalance: %+v\n", usrAddr.ToBase58(), testEnv.OntdBalance[usrAddr], testEnv.OnChainEState[0].ShareBalance[usrAddr], testEnv.OnChainEState[0].ShareBalance[usrAddr])

	tokenBought := big.NewInt(50)

	ontdBought1, tokenBought1 := testEnv.offTokenToTokenOutput(tokenBought)
	minOntdBought2, minTokenBought2 := ontdBought1, tokenBought1

	token1Hash, _ := common.AddressFromHexString(hex.EncodeToString(testEnv.OnChainTState[1].TokenAddr[:]))
	if err := testEnv.tokenToTokenOutput(0, tokenBought, minTokenBought2, minOntdBought2, testEnv.OnChainEState[0].Providers[userIndex], testEnv.Users[userIndex].Address, token1Hash); err != nil {
		log.Errorf("address: %s, ongToTokenSwapInput() error: %+v", err)
	}
	ontdBought3, tokenBought3 := testEnv.offTokenToTokenInput(tokenBought)
	minOntdBought4, minTokenBought4 := ontdBought3, tokenBought3
	if err := testEnv.tokenToTokenOutput(0, tokenBought, minTokenBought4, minOntdBought4, testEnv.OnChainEState[0].Providers[userIndex], testEnv.Users[1].Address, token1Hash); err != nil {
		log.Errorf("address: %s, ongToTokenSwapInput() error: %+v", err)
	}
}

func Test_tokenToExchangeInput(t *testing.T) {
	usrAddr := testEnv.Users[0].Address
	log.Infof("account: %s, ontdBalance: %+v, tokenBalance: %+v, shareBalance: %+v\n", usrAddr.ToBase58(), testEnv.OntdBalance[usrAddr], testEnv.OnChainEState[0].ShareBalance[usrAddr], testEnv.OnChainEState[0].ShareBalance[usrAddr])

	tokenSold := big.NewInt(20)

	ontdBought, tokenBought := testEnv.offTokenToTokenInput(tokenSold)
	minOntdBought, minTokenBought := ontdBought, tokenBought

	exchange1Hash, _ := common.AddressFromHexString(hex.EncodeToString(testEnv.OnChainEState[1].ExchangeAddr[:]))

	if err := testEnv.tokenToExchangeInput(0, tokenSold, minTokenBought, minOntdBought, testEnv.OnChainEState[0].Providers[0], testEnv.Users[0].Address, exchange1Hash); err != nil {
		log.Errorf("address: %s, ongToTokenSwapInput() error: %+v", err)
	}
	ontdBought, tokenBought = testEnv.offTokenToTokenInput(tokenSold)
	minOntdBought, minTokenBought = ontdBought, tokenBought
	if err := testEnv.tokenToExchangeInput(0, tokenSold, minTokenBought, minOntdBought, testEnv.OnChainEState[0].Providers[0], testEnv.Users[1].Address, exchange1Hash); err != nil {
		log.Errorf("address: %s, ongToTokenSwapInput() error: %+v", err)
	}
}

func Test_tokenToExchangeOutput(t *testing.T) {
	usrAddr := testEnv.Users[0].Address
	log.Infof("account: %s, ontdBalance: %+v, tokenBalance: %+v, shareBalance: %+v\n", usrAddr.ToBase58(), testEnv.OntdBalance[usrAddr], testEnv.OnChainEState[0].ShareBalance[usrAddr], testEnv.OnChainEState[0].ShareBalance[usrAddr])

	tokenSold := big.NewInt(50)

	ontdBought, tokenBought := testEnv.offTokenToTokenOutput(tokenSold)
	minOntdBought, minTokenBought := ontdBought, tokenBought

	exchange1Hash, _ := common.AddressFromHexString(hex.EncodeToString(testEnv.OnChainEState[1].ExchangeAddr[:]))
	if err := testEnv.tokenToExchangeOutput(0, tokenSold, minTokenBought, minOntdBought, testEnv.OnChainEState[0].Providers[0], testEnv.Users[0].Address, exchange1Hash); err != nil {
		log.Errorf("address: %s, ongToTokenSwapInput() error: %+v", err)
	}
	ontdBought, tokenBought = testEnv.offTokenToTokenInput(tokenSold)
	minOntdBought, minTokenBought = ontdBought, tokenBought
	if err := testEnv.tokenToExchangeOutput(0, tokenSold, minTokenBought, minOntdBought, testEnv.OnChainEState[0].Providers[0], testEnv.Users[1].Address, exchange1Hash); err != nil {
		log.Errorf("address: %s, ongToTokenSwapInput() error: %+v", err)
	}
}
