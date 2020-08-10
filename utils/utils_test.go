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

package utils

import (
	"fmt"
	"github.com/ontio/ontology/common"
	"github.com/skyinglyh1/uniswap_v1_test/config"
	"testing"
	"time"
)

func Test_CompileDeployContract(t *testing.T) {
	if err := config.DefConfig.Init("../config.json"); err != nil {
		fmt.Println("DefConfig.Init error:", err)
		return
	}
	_, _, err := GetSdkAndAccount(config.DefConfig.OntRpcAddress, config.DefConfig.WalletPath, config.DefConfig.AcctPwd)
	if err != nil {
		fmt.Printf("GetSdkAndAccount error: %v", err)
	}
	avmCode, err := CompileContract("../uniswap_v1_contracts/uniswap-v1/contracts/uniswap_exchange.py")
	if err != nil {
		fmt.Printf("Err: %v", err)
		return
	}
	fmt.Printf("avmCode is %x\n", avmCode)
	contractAddress := common.AddressFromVmCode(avmCode)
	fmt.Printf("ContractAddress is %s\n", contractAddress.ToHexString())
	//txHash, err := DeployContract(sdk, avmCode, acct)
	//if err != nil {
	//	fmt.Printf("DeployContract err: %v", err)
	//}
	//
	//if _, err = sdk.WaitForGenerateBlock(30*time.Second, 2); err != nil {
	//	fmt.Printf("Waiting err: %v", err)
	//	return
	//}
	//fmt.Printf("txHash: %s\n", txHash.ToHexString())

}
