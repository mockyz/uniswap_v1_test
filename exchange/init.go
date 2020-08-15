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
	"fmt"
	ontology_go_sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
	"github.com/skyinglyh1/uniswap_v1_test/config"
	"github.com/skyinglyh1/uniswap_v1_test/log"
	"github.com/skyinglyh1/uniswap_v1_test/utils"
	"math/big"
	"time"
)

var testEnv *TestEnv

type OnChainFactoryState struct {
	FactoryAddr             common.Address
	ExchangeHashToTokenAddr map[string]common.Address
	TokenHahsToExchangeAddr map[string]common.Address
	IdToTokenAddr           map[uint64]common.Address
}

type OnChainExchangeState struct {
	ExchangeAddr common.Address
	TokenAddr    common.Address
	FactoryAddr  common.Address
	Providers    []*ontology_go_sdk.Account
	OntdLiquid   *big.Int
	TokenLiquid  *big.Int
	ShareBalance map[common.Address]*big.Int
	ShareSupply  *big.Int
}

type OnChainTokenState struct {
	TokenAddr  common.Address
	Balances   map[common.Address]*big.Int
	Allowances map[common.Address]*big.Int
	Supply     *big.Int
}

type TestEnv struct {
	Sdk        *ontology_go_sdk.OntologySdk
	OntdAddr   common.Address
	Users      []*ontology_go_sdk.Account
	OtherUsers []common.Address

	OnChainFState *OnChainFactoryState
	OnChainTState []*OnChainTokenState
	OnChainEState []*OnChainExchangeState

	OffChainFState *OnChainFactoryState
	OffChainTState []*OnChainTokenState
	OffChainEState []*OnChainExchangeState

	OntdBalance   map[common.Address]*big.Int
	OntdAllowance map[common.Address]map[common.Address]*big.Int

	GasPrice      uint64
	GasLimit      uint64
	WaitTxTimeOut time.Duration
}

func init() {
	log.InitLog(1, log.Stdout)
	if err := config.DefConfig.Init("../config.json"); err != nil {
		log.Errorf("DefConfig.Init error:", err)
		return
	}
	sdk, accts, err := utils.GetSdkAndAccountNew()
	if err != nil {
		fmt.Printf("GetSdkAndAccount error: %v", err)
		return
	}
	factoryHash, err := common.AddressFromHexString(config.DefConfig.FactoryHash)
	if err != nil {
		log.Errorf("Exchange1Hash: %s, AddressFromHexString error: %v", config.DefConfig.FactoryHash, err)
		return
	}
	ontdHash, err := common.AddressFromHexString(config.DefConfig.OntdHash)
	if err != nil {
		log.Errorf("OntdHash: %s, AddressFromHexString error: %v", config.DefConfig.OntdHash, err)
		return
	}
	token1Hash, err := common.AddressFromHexString(config.DefConfig.Token1Hash)
	if err != nil {
		log.Errorf("Exchange1Hash: %s, AddressFromHexString error: %v", config.DefConfig.Token1Hash, err)
		return
	}
	exchange1Hash, err := common.AddressFromHexString(config.DefConfig.Exchange1Hash)
	if err != nil {
		log.Errorf("Exchange1Hash: %s, AddressFromHexString error: %v", config.DefConfig.Exchange1Hash, err)
		return
	}

	ofs := &OnChainFactoryState{
		FactoryAddr:             factoryHash,
		ExchangeHashToTokenAddr: make(map[string]common.Address),
		TokenHahsToExchangeAddr: make(map[string]common.Address),
		IdToTokenAddr:           make(map[uint64]common.Address),
	}
	ots1 := &OnChainTokenState{
		TokenAddr:  token1Hash,
		Balances:   make(map[common.Address]*big.Int),
		Allowances: make(map[common.Address]*big.Int),
	}

	oes1 := &OnChainExchangeState{
		ExchangeAddr: exchange1Hash,
		ShareBalance: make(map[common.Address]*big.Int),
	}
	if config.DefConfig.Token2Hash != "" {
		token2Hash, err := common.AddressFromHexString(config.DefConfig.Token2Hash)
		if err != nil {
			log.Errorf("Exchange1Hash: %s, AddressFromHexString error: %v", config.DefConfig.Token2Hash, err)
			return
		}
		exchange2Hash, err := common.AddressFromHexString(config.DefConfig.Exchange2Hash)
		if err != nil {
			log.Errorf("Exchange1Hash: %s, AddressFromHexString error: %v", config.DefConfig.Exchange2Hash, err)
			return
		}

		ots2 := &OnChainTokenState{
			TokenAddr:  token2Hash,
			Balances:   make(map[common.Address]*big.Int),
			Allowances: make(map[common.Address]*big.Int),
		}

		oes2 := &OnChainExchangeState{
			ExchangeAddr: exchange2Hash,
			ShareBalance: make(map[common.Address]*big.Int),
		}
		testEnv = &TestEnv{
			Sdk:            sdk,
			Users:          accts,
			OnChainFState:  ofs,
			OnChainTState:  []*OnChainTokenState{ots1, ots2},
			OnChainEState:  []*OnChainExchangeState{oes1, oes2},
			OffChainFState: ofs,
			OffChainTState: []*OnChainTokenState{ots1, ots2},
			OffChainEState: []*OnChainExchangeState{oes1, oes2},
			GasPrice:       config.DefConfig.GasPrice,
			GasLimit:       config.DefConfig.GasLimit,
			WaitTxTimeOut:  time.Duration(config.DefConfig.WaitTxTimeOut) * time.Second,
			OntdBalance:    make(map[common.Address]*big.Int),
		}
	} else {

		testEnv = &TestEnv{
			Sdk:            sdk,
			Users:          accts,
			OnChainFState:  ofs,
			OnChainTState:  []*OnChainTokenState{ots1},
			OnChainEState:  []*OnChainExchangeState{oes1},
			OffChainFState: ofs,
			OffChainTState: []*OnChainTokenState{ots1},
			OffChainEState: []*OnChainExchangeState{oes1},
			GasPrice:       config.DefConfig.GasPrice,
			GasLimit:       config.DefConfig.GasLimit,
			WaitTxTimeOut:  time.Duration(config.DefConfig.WaitTxTimeOut) * time.Second,
			OntdBalance:    make(map[common.Address]*big.Int),
		}
	}

	testEnv.OntdAddr = ontdHash
	testEnv.OntdAllowance = make(map[common.Address]map[common.Address]*big.Int)

	for _, otherUser := range config.DefConfig.OtherUsers {
		userAddr, _ := common.AddressFromBase58(otherUser)
		testEnv.OtherUsers = append(testEnv.OtherUsers, userAddr)
	}

	if err := testEnv.refreshFstate(); err != nil {
		log.Errorf("refreshFstate error: %v", err)
		return
	}
	if err := testEnv.refreshAcctBalance(); err != nil {
		log.Errorf("refreshAcctBalance error: %v", err)
		return
	}
}

func (this *TestEnv) refreshFstate() error {
	for i := 0; i < len(this.OnChainTState); i++ {
		tokenHash := common.ToArrayReverse(this.OnChainTState[i].TokenAddr[:])
		res, err := GetMethod(this.Sdk, this.OnChainFState.FactoryAddr, "getExchange", []interface{}{tokenHash})
		if err != nil {
			return fmt.Errorf("refreshFstate, getExchange err: %v", err)
		}
		exchange, err := common.AddressParseFromBytes(res)
		if err != nil {
			return fmt.Errorf("refreshFstate, getExchange err: %v", err)
		}
		if exchange != this.OnChainEState[i].ExchangeAddr {
			return fmt.Errorf("refreshFstate, getExchange(%x): %x != config.ExchangeAddr: %x", tokenHash, exchange[:], this.OnChainEState[i].ExchangeAddr[:])
		}
		this.OnChainFState.TokenHahsToExchangeAddr[hex.EncodeToString(tokenHash)] = exchange

		exchangeHash := common.ToArrayReverse(this.OnChainEState[i].ExchangeAddr[:])
		res1, err := GetMethod(this.Sdk, this.OnChainFState.FactoryAddr, "getToken", []interface{}{exchangeHash})
		if err != nil {
			return fmt.Errorf("refreshFstate, getToken err: %v", err)
		}
		tokenAddr, err := common.AddressParseFromBytes(res1)
		if err != nil {
			return fmt.Errorf("refreshFstate, getToken err: %v", err)
		}
		if tokenAddr != this.OnChainTState[i].TokenAddr {
			return fmt.Errorf("refreshFstate, getToken(%x): %x != config.TokenAddr: %x", exchangeHash, tokenAddr[:], this.OnChainTState[i].TokenAddr[:])
		}
		this.OnChainFState.ExchangeHashToTokenAddr[hex.EncodeToString(exchangeHash)] = exchange

	}

	return nil
}

//
//func (this *TestEnv) AddOnChainExchangeState() error {
//	for i, exchange := range this.OnChainEState {
//		preExeRes, err := this.Sdk.NeoVM.PreExecInvokeNeoVMContract(exchange.ExchangeAddr, []interface{}{"tokenAddress", []interface{}{}})
//		if err != nil {
//			return fmt.Errorf("AddOnChainExchangeState, Within exchange, tokenAddress(): error: %v", err)
//		}
//		tokenAddrBs, err := preExeRes.Result.ToByteArray()
//		if err != nil {
//			return fmt.Errorf("AddOnChainExchangeState, Within exchange, tokenAddress(), Result ToByteArray, error: %v", err)
//		}
//		tokenAddr, err := common.AddressParseFromBytes(tokenAddrBs)
//		if err != nil {
//			return fmt.Errorf("CheckTokenAndFactoryAddress, AddressParseFromBytes error: %v", err)
//		}
//		var configTokenHash string
//		if i == 0 {
//			configTokenHash = config.DefConfig.Token1Hash
//		} else if i == 1{
//			configTokenHash = config.DefConfig.Token2Hash
//		}
//		tokenHash, err := common.AddressFromHexString(configTokenHash)
//		if err != nil {
//			return fmt.Errorf("configExchangeHash: %s, AddressFromHexString error: %v", configTokenHash, err)
//		}
//		if tokenAddr != tokenHash {
//			log.Errorf("AddOnChainExchangeState, Exchange.tokenHash: %s not equal config.tokenHash: %s", tokenAddr.ToHexString(), tokenHash.ToHexString())
//		}
//		this.OnChainEState[i].TokenAddr = tokenAddr
//
//
//		preExeRes, err = this.Sdk.NeoVM.PreExecInvokeNeoVMContract(exchange.ExchangeAddr, []interface{}{"factoryAddress", []interface{}{}})
//		if err != nil {
//			return fmt.Errorf("AddOnChainExchangeState, Within exchange, factoryAddress(): error: %v", err)
//		}
//		factoryAddrBs, err := preExeRes.Result.ToByteArray()
//		if err != nil {
//			return fmt.Errorf("AddOnChainExchangeState, Within exchange, factoryAddress(), Result ToByteArray, error: %v", err)
//		}
//		factoryAddr, err := common.AddressParseFromBytes(factoryAddrBs)
//		if err != nil {
//			return fmt.Errorf("AddOnChainExchangeState, AddressParseFromBytes, factoryAddress error: %v", err)
//		}
//		factoryHash, err := common.AddressFromHexString(config.DefConfig.FactoryHash)
//		if err != nil {
//			return fmt.Errorf("configExchangeHash: %s, AddressFromHexString error: %v", configTokenHash, err)
//		}
//		if factoryAddr != factoryHash {
//			log.Errorf("AddOnChainExchangeState, Exchange.factoryhash: %s not equal config.factory: %s", factoryAddr.ToHexString(), factoryHash.ToHexString())
//		}
//		this.OnChainEState[i].FactoryAddr = factoryAddr
//
//
//		exOngBal, err := this.Sdk.Native.Ong.BalanceOf(exchange.ExchangeAddr)
//		if err != nil {
//			return fmt.Errorf("AddOnChainExchangeState, Get ong balance of exchange hash:%s error: %v", exchange.ExchangeAddr.ToHexString(), err)
//		}
//		this.OnChainEState[i].OngLiquid = big.NewInt(0).SetUint64(exOngBal)
//
//		preExeRes, err = this.Sdk.NeoVM.PreExecInvokeNeoVMContract(exchange.ExchangeAddr, []interface{}{"totalSupply", []interface{}{}})
//		if err != nil {
//			return fmt.Errorf("AddOnChainExchangeState, Within exchange, totalSupply(): error: %v", err)
//		}
//		shareSupplyBs, err := preExeRes.Result.ToByteArray()
//		if err != nil {
//			return fmt.Errorf("AddOnChainExchangeState, Within exchange, totalSupply(), Result ToByteArray, error: %v", err)
//		}
//		this.OnChainEState[i].ShareSupply = common.BigIntFromNeoBytes(shareSupplyBs)
//
//		token := this.OnChainEState[i]
//		exTokenBRes, err := this.Sdk.NeoVM.PreExecInvokeNeoVMContract(token.TokenAddr, []interface{}{"balanceOf", []interface{}{exchange.ExchangeAddr}})
//		if err != nil {
//			return fmt.Errorf("AddOnChainExchangeState, PreExecInvokeNeoVMContract, Get token:%s balance of exchange address:%s error: %v", token.TokenAddr.ToHexString(), exchange.ExchangeAddr.ToHexString(), err)
//		}
//		exTokenBs, err := exTokenBRes.Result.ToByteArray()
//		if err != nil {
//			return fmt.Errorf("AddOnChainExchangeState, PreExecute, Result.ToByteArray() Get token:%s balance of exchange address:%s, error: %v", token.TokenAddr.ToHexString(), exchange.ExchangeAddr.ToHexString(), err)
//		}
//		this.OnChainEState[i].TokenLiquid = common.BigIntFromNeoBytes(exTokenBs)
//
//
//		if err := this.refreshAcctBalance(); err != nil {
//			return fmt.Errorf("AddOnChainExchangeState, error: %v", err)
//		}
//	}
//
//}

func (this *TestEnv) refreshAcctBalance() error {
	userAddrs := make([]common.Address, 0)
	for _, user := range this.Users {
		userAddrs = append(userAddrs, user.Address)
	}
	for _, otherUser := range this.OtherUsers {
		userAddrs = append(userAddrs, otherUser)
	}
	//TODO: providers count control
	providers := this.Users

	exAddrs := make([]common.Address, 0)
	for _, exchange := range this.OnChainEState {
		exAddrs = append(exAddrs, exchange.ExchangeAddr)
	}

	tokenAddrs := make([]common.Address, 0)
	for _, v := range this.OnChainTState {
		tokenAddrs = append(tokenAddrs, v.TokenAddr)
	}

	// update TState at user balance, and user ong balance
	for _, userAddr := range userAddrs {
		balances, err := GetBalances(this.Sdk, userAddr, append(tokenAddrs, this.OntdAddr))
		if err != nil {
			return fmt.Errorf("refreshAcctBal, err: %v", err)
		}
		this.OntdBalance[userAddr] = balances[this.OntdAddr]
		for i, tokenAddr := range tokenAddrs {
			this.OnChainTState[i].Balances[userAddr] = balances[tokenAddr]
		}
		this.OntdAllowance[userAddr] = make(map[common.Address]*big.Int)

		allowances, err := GetAllowances(this.Sdk, this.OntdAddr, userAddr, exAddrs)
		if err != nil {
			return fmt.Errorf("refreshAcctBal, err: %v", err)
		}
		for spender, allowance := range allowances {
			this.OntdAllowance[userAddr][spender] = allowance
		}
	}

	//update TState at user allowance
	for _, exAddr := range exAddrs {
		for j, tokenAddr := range tokenAddrs {
			alls, supply, err := GetAllowancesAndSupply(this.Sdk, tokenAddr, exAddr, userAddrs)
			if err != nil {
				return fmt.Errorf("refreshAcctBal, err: %v", err)
			}
			for k, v := range alls {
				this.OnChainTState[j].Allowances[k] = v
			}
			this.OnChainTState[j].Supply = supply
		}
	}

	// update Estate at token liquid and ong liquid, tokenAddr and FactoryAddr
	for i, exAddr := range exAddrs {
		balances, err := GetBalances(this.Sdk, exAddr, append(tokenAddrs, this.OntdAddr))
		if err != nil {
			return fmt.Errorf("refreshAcctBal, err: %v", err)
		}
		// udpate ong liquid and token liquid
		this.OnChainEState[i].TokenLiquid = balances[tokenAddrs[i]]
		this.OnChainEState[i].OntdLiquid = balances[this.OntdAddr]

		// udpate token addr, factory addr
		ta, err := GetMethod(this.Sdk, this.OnChainEState[i].ExchangeAddr, "tokenAddress", nil)
		if err != nil {
			return fmt.Errorf("refershAcctBal, GetMethod, err: %v", err)
		}
		this.OnChainEState[i].TokenAddr, err = common.AddressParseFromBytes(ta)
		if err != nil {
			return fmt.Errorf("refershAcctBal, AddressParseFromBytes, err: %v", err)
		}

		fa, err := GetMethod(this.Sdk, this.OnChainEState[i].ExchangeAddr, "factoryAddress", nil)
		if err != nil {
			return fmt.Errorf("refershAcctBal, GetMethod, err: %v", err)
		}
		this.OnChainEState[i].FactoryAddr, err = common.AddressParseFromBytes(fa)
		if err != nil {
			return fmt.Errorf("refershAcctBal, AddressParseFromBytes, err: %v", err)
		}
		//	update providers, providers's shares and share supply
		this.OnChainEState[i].Providers = providers
		for _, provider := range this.OnChainEState[i].Providers {
			balances, err := GetBalances(this.Sdk, provider.Address, []common.Address{this.OnChainEState[i].ExchangeAddr})
			if err != nil {
				return fmt.Errorf("refreshAcctBal, err: %v", err)
			}
			this.OnChainEState[i].ShareBalance[provider.Address] = balances[this.OnChainEState[i].ExchangeAddr]
		}
		supplyRes, err := this.Sdk.NeoVM.PreExecInvokeNeoVMContract(this.OnChainEState[i].ExchangeAddr, []interface{}{"totalSupply", []interface{}{}})
		if err != nil {
			return fmt.Errorf("PreExec supply  error: %v", err)
		}
		supplyBs, err := supplyRes.Result.ToByteArray()
		if err != nil {
			return fmt.Errorf("Supply.ToByteArray()  error: %v", err)
		}
		this.OnChainEState[i].ShareSupply = common.BigIntFromNeoBytes(supplyBs)

	}

	return nil
}
func GetAllowances(sdk *ontology_go_sdk.OntologySdk, tokenAddr, owner common.Address, spenders []common.Address) (map[common.Address]*big.Int, error) {
	allowances := make(map[common.Address]*big.Int, 0)
	for _, spender := range spenders {
		allowance := new(big.Int)
		if tokenAddr == ontology_go_sdk.ONG_CONTRACT_ADDRESS {
			ontdBalance, err := sdk.Native.Ong.Allowance(owner, spender)
			if err != nil {
				return nil, fmt.Errorf("Get ong allowance of %s error: %v", owner.ToBase58(), err)
			}
			allowance = big.NewInt(0).SetUint64(ontdBalance)
		} else if tokenAddr == ontology_go_sdk.ONT_CONTRACT_ADDRESS {
			if tokenAddr == ontology_go_sdk.ONG_CONTRACT_ADDRESS {
				ontdBalance, err := sdk.Native.Ont.Allowance(owner, spender)
				if err != nil {
					return nil, fmt.Errorf("Get ong allowance of %s error: %v", owner.ToBase58(), err)
				}
				allowance = big.NewInt(0).SetUint64(ontdBalance)
			}
		} else {
			tokenBalanceRes, err := sdk.NeoVM.PreExecInvokeNeoVMContract(tokenAddr, []interface{}{"allowance", []interface{}{owner, spender}})
			if err != nil {
				return nil, fmt.Errorf("Get token:%s GetAllowances %s error: %v", tokenAddr.ToHexString(), owner.ToBase58(), err)
			}
			tokenBalanceBs, err := tokenBalanceRes.Result.ToByteArray()
			if err != nil {
				return nil, fmt.Errorf("Result.ToInteger() Get token:%s GetAllowances of %s error: %v", tokenAddr.ToHexString(), owner.ToBase58(), err)
			}
			allowance = common.BigIntFromNeoBytes(tokenBalanceBs)
		}
		allowances[spender] = allowance
	}
	return allowances, nil
}
func GetBalances(sdk *ontology_go_sdk.OntologySdk, owner common.Address, tokens []common.Address) (map[common.Address]*big.Int, error) {
	balances := make(map[common.Address]*big.Int, 0)
	for _, tokenAddr := range tokens {
		balance := new(big.Int)
		if tokenAddr == ontology_go_sdk.ONG_CONTRACT_ADDRESS {
			ontdBalance, err := sdk.Native.Ong.BalanceOf(owner)
			if err != nil {
				return nil, fmt.Errorf("Get ong balance of %s error: %v", owner.ToBase58(), err)
			}
			balance = big.NewInt(0).SetUint64(ontdBalance)
		} else if tokenAddr == ontology_go_sdk.ONT_CONTRACT_ADDRESS {
			if tokenAddr == ontology_go_sdk.ONG_CONTRACT_ADDRESS {
				ontdBalance, err := sdk.Native.Ont.BalanceOf(owner)
				if err != nil {
					return nil, fmt.Errorf("Get ong balance of %s error: %v", owner.ToBase58(), err)
				}
				balance = big.NewInt(0).SetUint64(ontdBalance)
			}
		} else {
			tokenBalanceRes, err := sdk.NeoVM.PreExecInvokeNeoVMContract(tokenAddr, []interface{}{"balanceOf", []interface{}{owner}})
			if err != nil {
				return nil, fmt.Errorf("Get token:%s balanceOf %s error: %v", tokenAddr.ToHexString(), owner.ToBase58(), err)
			}
			tokenBalanceBs, err := tokenBalanceRes.Result.ToByteArray()
			if err != nil {
				return nil, fmt.Errorf("Result.ToInteger() Get token:%s balance of %s error: %v", tokenAddr.ToHexString(), owner.ToBase58(), err)
			}
			balance = common.BigIntFromNeoBytes(tokenBalanceBs)
		}
		balances[tokenAddr] = balance
	}
	return balances, nil
}

func GetAllowancesAndSupply(sdk *ontology_go_sdk.OntologySdk, tokenAddr, spender common.Address, owners []common.Address) (map[common.Address]*big.Int, *big.Int, error) {
	allowances := make(map[common.Address]*big.Int, 0)
	for _, owner := range owners {
		allowRes, err := sdk.NeoVM.PreExecInvokeNeoVMContract(tokenAddr, []interface{}{"allowance", []interface{}{owner, spender}})
		if err != nil {
			return nil, nil, fmt.Errorf("Get token:%s balanceOf %s error: %v", tokenAddr.ToHexString(), owner.ToBase58(), err)
		}
		tokenBalanceBs, err := allowRes.Result.ToByteArray()
		if err != nil {
			return nil, nil, fmt.Errorf("Result.ToInteger() Get token:%s allowance(%s, %s) error: %v", tokenAddr.ToHexString(), owner.ToBase58(), spender.ToBase58(), err)
		}
		allowances[owner] = common.BigIntFromNeoBytes(tokenBalanceBs)
	}
	supplyRes, err := sdk.NeoVM.PreExecInvokeNeoVMContract(tokenAddr, []interface{}{"totalSupply", []interface{}{}})
	if err != nil {
		return nil, nil, fmt.Errorf("PreExec supply  error: %v", err)
	}
	supplyBs, err := supplyRes.Result.ToByteArray()
	if err != nil {
		return nil, nil, fmt.Errorf("Supply.ToByteArray()  error: %v", err)
	}
	return allowances, common.BigIntFromNeoBytes(supplyBs), nil
}

func GetMethod(sdk *ontology_go_sdk.OntologySdk, contractAddr common.Address, methodName string, params []interface{}) ([]byte, error) {
	if params == nil {
		params = []interface{}{}
	}
	preExeRes, err := sdk.NeoVM.PreExecInvokeNeoVMContract(contractAddr, []interface{}{methodName, params})
	if err != nil {
		return nil, fmt.Errorf("GetMethod, contractHash: %s, method: %s, pre invoke error %v", contractAddr.ToHexString(), methodName, err)
	}
	resBs, err := preExeRes.Result.ToByteArray()
	if err != nil {
		return nil, fmt.Errorf("GetMethod, contractHash: %s, method: %s, toBytearray error %v", contractAddr.ToHexString(), methodName, err)
	}
	return resBs, nil
}

//func (this *TestEnv) AddOffChainExchangeState() error {
//	preExeRes, err := this.Sdk.NeoVM.PreExecInvokeNeoVMContract(this.ExchangeHash, []interface{}{"tokenAddress", []interface{}{}})
//	if err != nil {
//		return fmt.Errorf("AddOnChainExchangeState, Within exchange, tokenAddress(): error: %v", err)
//	}
//	tokenAddrBs, err := preExeRes.Result.ToByteArray()
//	if err != nil {
//		return fmt.Errorf("AddOnChainExchangeState, Within exchange, tokenAddress(), Result ToByteArray, error: %v", err)
//	}
//	tokenAddr, err := common.AddressParseFromBytes(tokenAddrBs)
//	if err != nil {
//		return fmt.Errorf("CheckTokenAndFactoryAddress, AddressParseFromBytes error: %v", err)
//	}
//	if tokenAddr != this.TokenHash {
//		log.Errorf("AddOnChainExchangeState, Exchange.tokenHash: %s not equal tokenHash: %s", tokenAddr.ToHexString(), this.TokenHash.ToHexString())
//	}
//
//	preExeRes, err = this.Sdk.NeoVM.PreExecInvokeNeoVMContract(this.ExchangeHash, []interface{}{"factoryAddress", []interface{}{}})
//	if err != nil {
//		return fmt.Errorf("AddOnChainExchangeState, Within exchange, factoryAddress(): error: %v", err)
//	}
//	factoryAddrBs, err := preExeRes.Result.ToByteArray()
//	if err != nil {
//		return fmt.Errorf("AddOnChainExchangeState, Within exchange, factoryAddress(), Result ToByteArray, error: %v", err)
//	}
//	factoryAddr, err := common.AddressParseFromBytes(factoryAddrBs)
//	if err != nil {
//		return fmt.Errorf("AddOnChainExchangeState, AddressParseFromBytes, factoryAddress error: %v", err)
//	}
//	if factoryAddr != this.FactoryHash {
//		log.Errorf("AddOnChainExchangeState, Exchange.factoryhash: %s not equal factory: %s", factoryAddr.ToHexString(), this.FactoryHash.ToHexString())
//	}
//
//	exOngBal, err := this.Sdk.Native.Ong.BalanceOf(this.ExchangeHash)
//	if err != nil {
//		return fmt.Errorf("AddOnChainExchangeState, Get ong balance of exchange hash:%s error: %v", this.ExchangeHash.ToHexString(), err)
//	}
//	exOngB := big.NewInt(0).SetUint64(exOngBal)
//
//	preExeRes, err = this.Sdk.NeoVM.PreExecInvokeNeoVMContract(this.ExchangeHash, []interface{}{"totalSupply", []interface{}{}})
//	if err != nil {
//		return fmt.Errorf("AddOnChainExchangeState, Within exchange, totalSupply(): error: %v", err)
//	}
//	shareSupplyBs, err := preExeRes.Result.ToByteArray()
//	if err != nil {
//		return fmt.Errorf("AddOnChainExchangeState, Within exchange, totalSupply(), Result ToByteArray, error: %v", err)
//	}
//
//
//	exTokenBRes, err := this.Sdk.NeoVM.PreExecInvokeNeoVMContract(this.TokenHash, []interface{}{"balanceOf", []interface{}{this.ExchangeHash}})
//	if err != nil {
//		return fmt.Errorf("AddOnChainExchangeState, PreExecInvokeNeoVMContract, Get token:%s balance of exchange address:%x error: %v", this.TokenHash.ToHexString(), this.ExchangeHash[:], err)
//	}
//	exTokenBs, err := exTokenBRes.Result.ToByteArray()
//	if err != nil {
//		return fmt.Errorf("AddOnChainExchangeState, PreExecute, Result.ToByteArray() Get token:%s balance of exchange address:%xerror: %v", this.TokenHash.ToHexString(), this.ExchangeHash[:], err)
//	}
//
//	this.OffChainE1State = &OnChainExchangeState{
//		ExchangeAddr: this.ExchangeHash,
//		TokenAddr:    tokenAddr,
//		FactoryAddr:  factoryAddr,
//		OngLiquid:    exOngB,
//		TokenLiquid:   common.BigIntFromNeoBytes(exTokenBs),
//		ShareSupply:  common.BigIntFromNeoBytes(shareSupplyBs),
//		ShareBalance: make(map[common.Address]*big.Int),
//	}
//	if err := this.refreshAcctBalance(); err != nil {
//		return fmt.Errorf("AddOnChainExchangeState, error: %v", err)
//	}
//	for k, v := range this.BalanceMap.ShareBalance {
//		this.OffChainE1State.ShareBalance[k] = v
//	}
//	return nil
//}
func CheckContractExist(sdk *ontology_go_sdk.OntologySdk, contractAddr string) (bool, error) {
	dc, err := sdk.GetSmartContract(contractAddr)
	if err != nil {
		return false, fmt.Errorf("CheckContractExist, error: %v", err)
	}
	if dc.GetRawCode() != nil {
		return true, nil
	}
	return false, nil
}
