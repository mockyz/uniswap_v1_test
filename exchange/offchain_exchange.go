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
	"github.com/ontio/ontology/common"
	"github.com/skyinglyh1/uniswap_v1_test/log"
	"math/big"
)

type OffChainExchangeState struct {
	TokenAddr    common.Address
	FactoryAddr  common.Address
	OngLiquid    uint64
	TokenLiquid  *big.Int
	ShareBalance map[common.Address]*big.Int
	ShareSupply  *big.Int
}

func (this *OnChainExchangeState) getInputPrice(inputAmt *big.Int, inputReserve *big.Int, outputReserve *big.Int) *big.Int {
	//if (inputReserve.Cmp(big.NewInt(0)) < 1 && outputReserve.Cmp(big.NewInt(0)) < 1) {
	//	panic("_getInputPrice, assert error")
	//}
	inputAmtWithFee := big.NewInt(0).Mul(inputAmt, big.NewInt(9975))
	numerator := big.NewInt(0).Mul(inputAmtWithFee, outputReserve)
	denominator := big.NewInt(0).Add(big.NewInt(0).Mul(inputReserve, big.NewInt(10000)), inputAmtWithFee)
	return big.NewInt(0).Div(numerator, denominator)
}

func (this *OnChainExchangeState) getOutputPrice(outputAmt *big.Int, inputReserve *big.Int, outputReserve *big.Int) *big.Int {
	numerator := big.NewInt(0).Mul(big.NewInt(0).Mul(inputReserve, outputAmt), big.NewInt(10000))
	denominator := big.NewInt(0).Mul(big.NewInt(0).Sub(outputReserve, outputAmt), big.NewInt(9975))
	return big.NewInt(0).Div(big.NewInt(0).Sub(big.NewInt(0).Add(numerator, denominator), big.NewInt(1)), denominator)
}

func (this *OnChainExchangeState) offOntToTokenInput(ontdSold *big.Int, minTokens *big.Int) error {
	tokenBought := this.getInputPrice(ontdSold, this.OntdLiquid, this.TokenLiquid)
	//TODO: exchange token balance increase ongBought
	//TODO: exchange ong decrease ongBought
	log.Debugf("offOntToTokenInput, ontToTokenInput, tokenBought is %+v, minTokens is %+v", tokenBought.String(), minTokens.String())
	return nil
}
func (this *OnChainExchangeState) offOntToTokenOutput(tokensBought *big.Int, maxOntd *big.Int) error {
	ontdSold := this.getOutputPrice(tokensBought, this.OntdLiquid, this.TokenLiquid)
	//TODO: exchange token balance increase ongBought
	//TODO: exchange ong decrease ongBought
	log.Debugf("offOntToTokenOutput, OntToTokenOutput, ontdSold is %+v, maxOntd is %+v", ontdSold.String(), maxOntd.String())
	return nil
}
func (this *OnChainExchangeState) offTokenToOntInput(tokenSold *big.Int, minOng *big.Int) error {
	ongBought := this.getInputPrice(tokenSold, this.TokenLiquid, this.OntdLiquid)
	//TODO: exchange token balance increase ongBought
	//TODO: exchange ong decrease ongBought
	log.Debugf("offTokenToOntInput, ongBought is %+v, minOng is %+v", ongBought.String(), minOng.String())
	return nil
}
func (this *OnChainExchangeState) offTokenToOntOutput(ongBought *big.Int, maxTokens *big.Int) error {
	tokenSold := this.getOutputPrice(ongBought, this.TokenLiquid, this.OntdLiquid)
	//TODO: exchange token balance increase ongBought
	//TODO: exchange ong decrease ongBought
	log.Debugf("offTokenToOntOutput, tokenSold is %+v, maxToken is %+v", tokenSold.String(), maxTokens.String())
	return nil
}

func (this *TestEnv) offTokenToTokenInput(tokenSold *big.Int) (*big.Int, *big.Int) {
	ontdBought := this.OnChainEState[0].getInputPrice(tokenSold, this.OnChainEState[0].TokenLiquid, this.OnChainEState[0].OntdLiquid)
	//TODO: exchange token balance increase ontDBought
	//TODO: exchange ong decrease ontDBought
	tokenBought := this.OnChainEState[1].getInputPrice(ontdBought, this.OnChainEState[1].OntdLiquid, this.OnChainEState[1].TokenLiquid)
	log.Debugf("offTokenToTokenInput, tokenSold is %+v, minOntdBought is %+v,  minTokenBought is is %+v", tokenSold.String(), ontdBought, tokenBought.String())
	return ontdBought, tokenBought
}

func (this *TestEnv) offTokenToTokenOutput(tokenBought *big.Int) (*big.Int, *big.Int) {
	ontdBought := this.OnChainEState[0].getOutputPrice(tokenBought, this.OnChainEState[0].TokenLiquid, this.OnChainEState[0].OntdLiquid)
	//TODO: exchange token balance increase ongBought
	//TODO: exchange ong decrease ongBought
	tokenBought1 := this.OnChainEState[1].getInputPrice(ontdBought, this.OnChainEState[1].OntdLiquid, this.OnChainEState[1].TokenLiquid)
	log.Debugf("offTokenToTokenInput, tokenSold is %+v, minOntdBought is %+v,  minTokenBought is is %+v", tokenBought.String(), ontdBought, tokenBought.String())
	return ontdBought, tokenBought1
}
