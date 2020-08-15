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
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ontio/ontology-crypto/keypair"
	ontology_go_sdk "github.com/ontio/ontology-go-sdk"
	sdkcommon "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/types"
	"github.com/skyinglyh1/uniswap_v1_test/config"
	"github.com/skyinglyh1/uniswap_v1_test/log"
	"io/ioutil"
	"net/http"
	"strings"
)

type CompilePayLoad struct {
	Type string `json:"type"`
	Code string `json:"code"`
}
type CompileResponse struct {
	ErrorCode uint64 `json:"errcode"`
	Avm       string `json:"avm"`
	Abi       string `json:"abi"`
	Debug     string `json:"debug"`
	Opcode    string `json:"opcode"`
	FuncMap   string `json:"funcmap"`
}

const CompilerUrl = "http://42.159.92.140:8089/api/v2.0/python/compile"

func GetSdkAndAccount(url, walletPath, passwd string) (*ontology_go_sdk.OntologySdk, []*ontology_go_sdk.Account, error) {
	err := config.DefConfig.Init(walletPath)
	if err != nil {
		fmt.Println("DefConfig.Init error:", err)
		return nil, nil, fmt.Errorf("DefConfig.Init error: %v", err)
	}
	ontSdk := ontology_go_sdk.NewOntologySdk()
	ontSdk.NewRpcClient().SetAddress(url)

	wallet, err := ontSdk.OpenWallet(walletPath)
	if err != nil {
		return nil, nil, fmt.Errorf("OpenWallet error: %v", err)
	}
	accts := make([]*ontology_go_sdk.Account, 0)
	count := wallet.GetAccountCount()
	for i := 1; i <= count; i++ {
		acct, err := wallet.GetAccountByIndex(i, []byte(passwd))
		if err != nil {
			return nil, nil, fmt.Errorf("wallet.GetDefaultAccount error: %v", err)
		}
		accts = append(accts, acct)
	}

	return ontSdk, accts, nil
}

//sdk, accts, err := utils.GetSdkAndAccount(config.DefConfig.OntRpcAddress, config.DefConfig.WalletPath, config.DefConfig.AcctPwd)

func GetSdkAndAccountNew() (*ontology_go_sdk.OntologySdk, []*ontology_go_sdk.Account, error) {
	ontSdk := ontology_go_sdk.NewOntologySdk()
	ontSdk.NewRpcClient().SetAddress(config.DefConfig.OntRpcAddress)
	accts := make([]*ontology_go_sdk.Account, 0)
	Wif := config.DefConfig.WIF
	for i := 0; i < len(Wif); i++ {
		acct, err := NewAccountByWif(Wif[i])
		if err != nil {
			return nil, nil, fmt.Errorf("wallet.GetDefaultAccount error: %v", err)
		}
		accts = append(accts, acct)
	}
	wallet, err := ontSdk.OpenWallet(config.DefConfig.WalletPath)
	if err != nil {
		return nil, nil, fmt.Errorf("OpenWallet error: %v", err)
	}
	count := wallet.GetAccountCount()
	for i := 1; i <= count; i++ {
		acct, err := wallet.GetAccountByIndex(i, []byte(config.DefConfig.AcctPwd))
		if err != nil {
			return nil, nil, fmt.Errorf("wallet.GetDefaultAccount error: %v", err)
		}
		accts = append(accts, acct)
	}
	return ontSdk, accts, nil
}
func CompileContract(contractFilePath string) ([]byte, error) {
	//contractsPath := config.ContractsPath
	//factoryFile := contractsPath + "uniswap_factory.py"
	//exchangeFile := contractsPath + "uniswap_exchange.py"
	data, err := ioutil.ReadFile(contractFilePath)
	if err != nil {
		return nil, fmt.Errorf("ReadFile: %s, err: %v", contractFilePath, err)
	}
	payload := CompilePayLoad{
		Type: "Python",
		Code: string(data),
	}
	payloadBs, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("jsonMarshal payload err: %v", err)
	}
	resp, err := http.Post(CompilerUrl, "application/json", bytes.NewReader(payloadBs))
	if err != nil {
		return nil, fmt.Errorf("http.Post error: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read rpc response body error:%s", err)
	}
	rpcRsp := new(CompileResponse)
	err = json.Unmarshal(body, rpcRsp)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal JsonRpcResponse:%s error:%s", body, err)
	}
	if rpcRsp.ErrorCode != 0 {
		return nil, fmt.Errorf("err code: %v", rpcRsp.ErrorCode)
	}
	avm := strings.Split(rpcRsp.Avm, "'")
	avmCode, err := hex.DecodeString(avm[1])
	if err != nil {
		return nil, fmt.Errorf("decode response.Avm err: %v", err)
	}
	return avmCode, nil
}

func CheckContractDeployed(sdk *ontology_go_sdk.OntologySdk, contractHash common.Address) bool {
	sdk.GetSmartContract(contractHash.ToHexString())
	return true
}

func DeployContract(sdk *ontology_go_sdk.OntologySdk, avmCode []byte, signer *ontology_go_sdk.Account) (*common.Uint256, error) {
	txHash, err := sdk.NeoVM.DeployNeoVMSmartContract(config.DefConfig.GasPrice, config.DefConfig.GasLimit*100000, signer, true, hex.EncodeToString(avmCode), "name", "Version", "author", "email", "desc")
	if err != nil {
		return nil, fmt.Errorf("DepolyContract, error: %v", err)
	}
	return &txHash, nil
}

func CheckContracts(sdk *ontology_go_sdk.OntologySdk, factoryPath, tokenPath string, factoryHash, tokenHash common.Address, filePriorHash bool) ([]common.Address, error) {
	newConHashes := make([]common.Address, 0)
	if filePriorHash {
		// Need to compile contract
		avmCode, err := CompileContract(factoryPath)
		if err != nil {
			return nil, fmt.Errorf("Compile contract with path %s error: %v", factoryPath, err)
		}
		newConH := common.AddressFromVmCode(avmCode)
		newConHashes = append(newConHashes, newConH)

		avmCode, err = CompileContract(factoryPath)
		if err != nil {
			return nil, fmt.Errorf("Compile contract with path %s error: %v", tokenPath, err)
		}
		newConH = common.AddressFromVmCode(avmCode)
		newConHashes = append(newConHashes, newConH)

	}
	//for _, file := range hashToContractFiles {
	//	avmCode, err := CompileContract(file)
	//	if err != nil {
	//		return fmt.Errorf("Compile contract with path %s error: %v", )
	//	}
	//	contractAddr := common.AddressFromVmCode(avmCode)
	//
	//}
	return nil, nil
}

func PrintSmartEventByHash_Ont(sdk *ontology_go_sdk.OntologySdk, txHash string) []*sdkcommon.NotifyEventInfo {
	evts, err := sdk.GetSmartContractEvent(txHash)
	if err != nil {
		fmt.Printf("GetSmartContractEvent error:%s", err)
		return nil
	}
	fmt.Printf("evts = %+v\n", evts)
	fmt.Printf("TxHash:%s\n", txHash)
	fmt.Printf("State:%d\n", evts.State)
	for _, notify := range evts.Notify {
		fmt.Printf("ContractAddress:%s\n", notify.ContractAddress)
		fmt.Printf("States:%+v\n", notify.States)
	}
	return evts.Notify
}

func NewAccountByWif(Wif string) (*ontology_go_sdk.Account, error) {
	// AScExXzLbkZV32tDFdV7Uoq7ZhCT1bRCGp
	privateKey, err := keypair.WIF2Key([]byte(Wif))
	if err != nil {
		log.Errorf("decrypt privateKey error:%s", err)
	}
	pub := privateKey.Public()
	address := types.AddressFromPubKey(pub)
	log.Infof("address: %s\n", address.ToBase58())
	return &ontology_go_sdk.Account{
		PrivateKey: privateKey,
		PublicKey:  pub,
		Address:    address,
	}, nil
}
