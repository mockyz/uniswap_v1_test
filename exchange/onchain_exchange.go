package exchange

import (
	"fmt"
	ontology_go_sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
	"github.com/skyinglyh1/uniswap_v1_test/log"
	"github.com/skyinglyh1/uniswap_v1_test/utils"
	"math/big"
	"time"
)


type ExchangeProviderState struct {
	ShareBalance map[common.Address]*big.Int
}






type ProviderBalanceMap struct {
	OngBalance map[common.Address]uint64
	TokenBalance map[common.Address]*big.Int
	AllowanceBalance map[common.Address]*big.Int
	ShareBalance map[common.Address]*big.Int

}



func (this *TestEnv) addLiquid(exchangeIndex int, minLiquidity, maxTokens *big.Int, ontdAmt *big.Int) error {
	if err := this.refreshAcctBalance(); err != nil {
		return fmt.Errorf("addLiquid, refreshBalance err: %v", err)
	}

	exOntdBalance1 := this.OnChainEState[exchangeIndex].OntdLiquid
	exTokenBalance1 := this.OnChainEState[exchangeIndex].TokenLiquid

	for _, provider := range this.OnChainEState[exchangeIndex].Providers {
		if this.OnChainTState[0].Balances[provider.Address].Cmp(maxTokens) < 0 {
			return fmt.Errorf("provider: %s does not have enough token: %v", provider.Address.ToBase58(), maxTokens)
		}
		if this.OnChainTState[0].Allowances[provider.Address].Cmp(maxTokens) < 0 {
			// approve token to exchange
			approveTxHash, err := this.Sdk.NeoVM.InvokeNeoVMContract(this.GasPrice, this.GasLimit, provider, provider, this.OnChainEState[exchangeIndex].TokenAddr, []interface{}{"approve", []interface{}{
				provider.Address, this.OnChainEState[exchangeIndex].ExchangeAddr, maxTokens,
			}})
			if err != nil {
				return fmt.Errorf("Provider: %s, approve token to exchange err: %v", provider.Address.ToBase58(), err)
			}
			if _, err := this.Sdk.WaitForGenerateBlock(this.WaitTxTimeOut, 1); err != nil {
				return fmt.Errorf("Ontology, not generate block after %+v, err: %v", this.WaitTxTimeOut, err)
			}
			utils.PrintSmartEventByHash_Ont(this.Sdk, approveTxHash.ToHexString())

		}
		if this.OntdBalance[provider.Address].Cmp(ontdAmt) < 0 {
			return fmt.Errorf("provider: %s does not have enough ontd: %v", provider.Address.ToBase58(), ontdAmt)
		}

		// addLiquidity
		txHash, err := this.Sdk.NeoVM.InvokeNeoVMContract(this.GasPrice, this.GasLimit, provider, provider, this.OnChainEState[exchangeIndex].ExchangeAddr, []interface{}{"addLiquidity", []interface{}{
			minLiquidity,
				//TODO: check
				maxTokens,
				time.Now().Add(this.WaitTxTimeOut).Unix(), provider.Address, ontdAmt,
		}})
		if err != nil {
			return fmt.Errorf("Provider: %s, addLiquid err: %v", provider.Address.ToBase58(), err)
		}
		if _, err := this.Sdk.WaitForGenerateBlock(this.WaitTxTimeOut, 1); err != nil {
			return fmt.Errorf("Ontology, not generate block after %+v, err: %v", this.WaitTxTimeOut, err)
		}
		utils.PrintSmartEventByHash_Ont(this.Sdk, txHash.ToHexString())
	}


	if err := this.refreshAcctBalance(); err != nil {
		return fmt.Errorf("addLiquid, refreshBalance err: %v", err)
	}
	exOngBalance2 := this.OnChainEState[exchangeIndex].OntdLiquid
	exTokenBalance2 := this.OnChainEState[exchangeIndex].TokenLiquid

	if big.NewInt(0).Sub( exOngBalance2, exOntdBalance1).Cmp(ontdAmt) < 0{
		// TODO: if we don't count the tx fee, they will not equal
		return fmt.Errorf("exchange ong balance increse incorrect")
	}
	//if big.NewInt(0).Sub(exTokenBalance2, exTokenBalance1).Cmp(maxTokens) != 1 {
	//	return fmt.Errorf("exchange token balance increse incorrect")
	//}
	log.Debugf("exchange token increased by : %v\n", big.NewInt(0).Sub(exTokenBalance2, exTokenBalance1))

	// TODO: update off chain state
	// TODO: Check onchain states equals offchain states

	return nil
}


func (this *TestEnv) removeLiquid(exchangeIndex int, amount, min_ong *big.Int, withdrawer *ontology_go_sdk.Account) error {
	if err := this.refreshAcctBalance(); err != nil {
		return fmt.Errorf("removeLiquid, refreshBalance err: %v", err)
	}

	exOngBalance1 := this.OnChainEState[exchangeIndex].OntdLiquid
	shareB1 := this.OnChainEState[exchangeIndex].ShareBalance[withdrawer.Address]

	// Condition check
	if this.OnChainEState[exchangeIndex].ShareBalance[withdrawer.Address].Cmp(amount) < 0 {
		return fmt.Errorf("removeLiquid, withdrawer: %s, not have enough share balance", withdrawer.Address.ToBase58())
	}

	// removeLiquidity
	txHash, err := this.Sdk.NeoVM.InvokeNeoVMContract(this.GasPrice, this.GasLimit, withdrawer, withdrawer, this.OnChainEState[exchangeIndex].ExchangeAddr, []interface{}{"removeLiquidity", []interface{}{
		amount,
		1,
		1,
		time.Now().Add(this.WaitTxTimeOut).Unix(),
		withdrawer.Address,
	}})
	if err != nil {
		return fmt.Errorf("removeLiquid, withdrawer: %s withdraw err: %v", withdrawer.Address.ToBase58(), err)
	}
	if _, err := this.Sdk.WaitForGenerateBlock(this.WaitTxTimeOut, 1); err != nil {
		return fmt.Errorf("removeLiquid, Ontology, not generate block after %+v, err: %v", this.WaitTxTimeOut, err)
	}
	utils.PrintSmartEventByHash_Ont(this.Sdk, txHash.ToHexString())



	if err := this.refreshAcctBalance(); err != nil {
		return fmt.Errorf("removeLiquid, refreshBalance err: %v", err)
	}
	exOngBalance2 := this.OnChainEState[exchangeIndex].OntdLiquid
	shareB2 := this.OnChainEState[exchangeIndex].ShareBalance[withdrawer.Address]

	//if exOngBalance2 - exOngBalance1 != min_ong {
	//	return fmt.Errorf("exchange ong balance decrease incorrect")
	//}
	log.Debugf("removeLiquid, exchange ong decreased by : %d\n", exOngBalance1.Uint64() - exOngBalance2.Uint64())
	// TODO: token means share balance
	if big.NewInt(0).Sub(shareB1, shareB2).Cmp(amount) != 0 {
		return fmt.Errorf("removeLiquid, withdrawer share balance decrease incorrect")
	}

	return nil
}



func (this *TestEnv) ontToTokenInput(ontdAmt, minTokens *big.Int, invoker *ontology_go_sdk.Account, recipient common.Address) error {
	if err := this.refreshAcctBalance(); err != nil {
		return fmt.Errorf("ontToTokenInput, refreshBalance err: %v", err)
	}

	exOngBalance1 := this.OnChainEState[0].OntdLiquid
	recBal1 := this.OnChainTState[0].Balances[recipient]
	// Condition check
	if this.OntdBalance[invoker.Address].Cmp(ontdAmt) < 0 {
		return fmt.Errorf("ontToTokenInput, invoker: %s, not have enough ontd balance", invoker.Address.ToBase58())
	}

	if this.OntdAllowance[invoker.Address][this.OnChainEState[0].ExchangeAddr].Cmp(ontdAmt) < 0 {
		// approve ontd  to exchange
		approveTxHash, err := this.Sdk.NeoVM.InvokeNeoVMContract(this.GasPrice, this.GasLimit, invoker, invoker, this.OntdAddr, []interface{}{"approve", []interface{}{
			invoker.Address, this.OnChainEState[0].ExchangeAddr, ontdAmt,
		}})
		if err != nil {
			return fmt.Errorf("ontToTokenInput: %s, approve token to exchange err: %v", invoker.Address.ToBase58(), err)
		}
		if _, err := this.Sdk.WaitForGenerateBlock(this.WaitTxTimeOut, 1); err != nil {
			return fmt.Errorf("Ontology, not generate block after %+v, err: %v", this.WaitTxTimeOut, err)
		}
		utils.PrintSmartEventByHash_Ont(this.Sdk, approveTxHash.ToHexString())

	}

	var params []interface{}
	// ongToTokenSwap
	if invoker.Address == recipient {
		params = []interface{}{
			"ontToTokenSwapInput",
			[]interface{}{
				minTokens,
				time.Now().Add(this.WaitTxTimeOut).Unix(),
				invoker.Address,
				ontdAmt,
			},
		}
	} else {
		params = []interface{}{
			"ontToTokenTransferInput",
			[]interface{}{
				minTokens,
				time.Now().Add(this.WaitTxTimeOut).Unix(),
				recipient,
				invoker.Address,
				ontdAmt,
			},
		}
	}
	txHash, err := this.Sdk.NeoVM.InvokeNeoVMContract(this.GasPrice, this.GasLimit, invoker, invoker, this.OnChainEState[0].ExchangeAddr, params)
	if err != nil {
		return fmt.Errorf("ontToTokenInput, invoker: %s invoke err: %v", invoker.Address.ToBase58(), err)
	}
	if _, err := this.Sdk.WaitForGenerateBlock(this.WaitTxTimeOut, 1); err != nil {
		return fmt.Errorf("Ontology, not generate block after %+v, err: %v", this.WaitTxTimeOut, err)
	}
	utils.PrintSmartEventByHash_Ont(this.Sdk, txHash.ToHexString())



	if err := this.refreshAcctBalance(); err != nil {
		return fmt.Errorf("ontToTokenInput, refreshBalance err: %v", err)
	}
	exOngBalance2 := this.OnChainEState[0].OntdLiquid

	ongDecrement := big.NewInt(0).Sub(exOngBalance2, exOngBalance1)
	if ongDecrement.Cmp(ontdAmt) < 0 {
		return fmt.Errorf("exchange ong balance decrease incorrect")
	}
	log.Debugf("exchange ong increased by : %d\n", ongDecrement)

	if invoker.Address == recipient {

	} else {
		recBal2 := this.OnChainTState[0].Balances[recipient]
		increment := big.NewInt(0).Sub(recBal2, recBal1)
		if increment.Cmp(big.NewInt(0)) > 1 {
			log.Debugf("recipient: %s received %+v tokens\n", recipient.ToBase58(), increment)
		}
	}

	return nil
}



func (this *TestEnv) ontToTokenOutput(tokenBought *big.Int, maxOntd *big.Int, invoker *ontology_go_sdk.Account, recipient common.Address) error {
	if err := this.refreshAcctBalance(); err != nil {
		return fmt.Errorf("ongToTokenSwapInput, refreshBalance err: %v", err)
	}

	exOngBalance1 := this.OnChainEState[0].OntdLiquid
	exTokenB1 := this.OnChainEState[0].TokenLiquid
	recBal1 := this.OnChainTState[0].Balances[recipient]
	// Condition check
	if this.OntdBalance[invoker.Address].Cmp(maxOntd) < 0 {
		return fmt.Errorf("ongToTokenOutput, invoker: %s, not have enough ong balance", invoker.Address.ToBase58())
	}

	if this.OntdAllowance[invoker.Address][this.OnChainEState[0].ExchangeAddr].Cmp(maxOntd) < 0 {
		// approve ontd  to exchange
		approveTxHash, err := this.Sdk.NeoVM.InvokeNeoVMContract(this.GasPrice, this.GasLimit, invoker, invoker, this.OntdAddr, []interface{}{"approve", []interface{}{
			invoker.Address, this.OnChainEState[0].ExchangeAddr, maxOntd,
		}})
		if err != nil {
			return fmt.Errorf("ontToTokenInput: %s, approve token to exchange err: %v", invoker.Address.ToBase58(), err)
		}
		if _, err := this.Sdk.WaitForGenerateBlock(this.WaitTxTimeOut, 1); err != nil {
			return fmt.Errorf("Ontology, not generate block after %+v, err: %v", this.WaitTxTimeOut, err)
		}
		utils.PrintSmartEventByHash_Ont(this.Sdk, approveTxHash.ToHexString())

	}
	var params []interface{}
	// ongToTokenSwap
	if invoker.Address == recipient {
		params = []interface{}{
			"ontToTokenSwapOutput",
			[]interface{}{
				tokenBought,
				time.Now().Add(this.WaitTxTimeOut).Unix(),
				invoker.Address,
				maxOntd,
			},
		}
	} else {
		params = []interface{}{
			"ontToTokenTransferOutput",
			[]interface{}{
				tokenBought,
				time.Now().Add(this.WaitTxTimeOut).Unix(),
				recipient,
				invoker.Address,
				maxOntd,
			},
		}
	}
	txHash, err := this.Sdk.NeoVM.InvokeNeoVMContract(this.GasPrice, this.GasLimit, invoker, invoker, this.OnChainEState[0].ExchangeAddr, params)
	if err != nil {
		return fmt.Errorf("ongToTokenOutput, invoker: %s invoke err: %v", invoker.Address.ToBase58(), err)
	}
	if _, err := this.Sdk.WaitForGenerateBlock(this.WaitTxTimeOut, 1); err != nil {
		return fmt.Errorf("Ontology, not generate block after %+v, err: %v", this.WaitTxTimeOut, err)
	}
	utils.PrintSmartEventByHash_Ont(this.Sdk, txHash.ToHexString())


	if err := this.refreshAcctBalance(); err != nil {
		return fmt.Errorf("ongToTokenOutput, refreshBalance err: %v", err)
	}
	exOngBalance2 := this.OnChainEState[0].OntdLiquid
	exTokenB2 := this.OnChainEState[0].TokenLiquid
	ongIncrement := big.NewInt(0).Sub(exOngBalance2, exOngBalance1)
	if ongIncrement.Cmp(maxOntd) > 0 {
		return fmt.Errorf("ongToTokenOutput, exchange ong balance increase incorrect")
	}
	log.Debugf("ongToTokenOutput, exchange ong increased by : %d\n", ongIncrement.String())
	log.Debugf("ongToTokenOutput, exchange token decreased by : %d\n", exTokenB1.Sub(exTokenB1, exTokenB2).String())

	if invoker.Address == recipient {

	} else {
		recBal2 := this.OnChainTState[0].Balances[recipient]
		increment := big.NewInt(0).Sub(recBal2, recBal1)
		if increment.Cmp(big.NewInt(0)) > 1 {
			log.Debugf("ongToTokenOutput, recipient: %s received %+v tokens\n", recipient.ToBase58(), increment)
		}
	}

	return nil
}



func (this *TestEnv) tokenToOntInput(tokenSold *big.Int, minOng *big.Int, invoker *ontology_go_sdk.Account, recipient common.Address) error {
	if err := this.refreshAcctBalance(); err != nil {
		return fmt.Errorf("tokenToOngInput, refreshBalance err: %v", err)
	}

	exOngBalance1 := this.OnChainEState[0].OntdLiquid
	exTokenB1 := this.OnChainEState[0].TokenLiquid
	recOngB1 := this.OnChainTState[0].Balances[recipient]

	// Condition check
	if this.OnChainEState[0].TokenLiquid.Cmp(tokenSold) < 0 {
		return fmt.Errorf("tokenToOngInput, exchange token balance: %v < tokenSold: %v", this.OnChainEState[0].TokenLiquid.String(), tokenSold.String())
	}
	if this.OnChainTState[0].Balances[invoker.Address].Cmp(tokenSold) < 1 {
		return fmt.Errorf("tokenToOngInput, invoker: %s, not have enough token balance", invoker.Address.ToBase58())
	}
	if this.OnChainTState[0].Allowances[invoker.Address].Cmp(tokenSold) < 1 {
		log.Debugf("tokenToOngInput, invoker: %s, not have allowance token", invoker.Address.ToBase58())
		// approve token to exchange from invoker
		approveTxHash, err := this.Sdk.NeoVM.InvokeNeoVMContract(this.GasPrice, this.GasLimit, invoker, invoker, this.OnChainTState[0].TokenAddr, []interface{}{"approve", []interface{}{
			invoker.Address, this.OnChainEState[0].ExchangeAddr, tokenSold,
		}})
		if err != nil {
			return fmt.Errorf("tokenToOngInput: %s, approve token to exchange err: %v", invoker.Address.ToBase58(), err)
		}
		if _, err := this.Sdk.WaitForGenerateBlock(this.WaitTxTimeOut, 1); err != nil {
			return fmt.Errorf("Ontology, not generate block after %+v, err: %v", this.WaitTxTimeOut, err)
		}
		utils.PrintSmartEventByHash_Ont(this.Sdk, approveTxHash.ToHexString())
	}

	var params []interface{}
	// tokenToOngInput
	if invoker.Address == recipient {
		params = []interface{}{
			"tokenToOntSwapInput",
			[]interface{}{
				tokenSold,
				minOng.Uint64(),
				time.Now().Add(this.WaitTxTimeOut).Unix(),
				invoker.Address,
			},
		}
	} else {
		params = []interface{}{
			"tokenToOntTransferInput",
			[]interface{}{
				tokenSold,
				minOng.Uint64(),
				time.Now().Add(this.WaitTxTimeOut).Unix(),
				invoker.Address,
				recipient,
			},
		}
	}
	txHash, err := this.Sdk.NeoVM.InvokeNeoVMContract(this.GasPrice, this.GasLimit, invoker, invoker, this.OnChainEState[0].ExchangeAddr, params)
	if err != nil {
		return fmt.Errorf("tokenToOngInput, invoker: %s invoke err: %v", invoker.Address.ToBase58(), err)
	}
	if _, err := this.Sdk.WaitForGenerateBlock(this.WaitTxTimeOut, 1); err != nil {
		return fmt.Errorf("Ontology, not generate block after %+v, err: %v", this.WaitTxTimeOut, err)
	}
	utils.PrintSmartEventByHash_Ont(this.Sdk, txHash.ToHexString())



	if err := this.refreshAcctBalance(); err != nil {
		return fmt.Errorf("tokenToOngInput, refreshBalance err: %v", err)
	}
	exOngBalance2 := this.OnChainEState[0].OntdLiquid
	exTokenB2 := this.OnChainEState[0].TokenLiquid
	exOngInc := big.NewInt(0).Sub(exOngBalance2, exOngBalance1)
	exTokenInc := big.NewInt(0).Sub(exTokenB2, exTokenB1)
	log.Debugf("exchange ong increased by : %s\n", exOngInc.String())
	log.Debugf("exchange token increased by : %s\n", exTokenInc.String())

	if invoker.Address == recipient {

	} else {
		recOngB2 := this.OnChainTState[0].Balances[recipient]
		ongIncrement := big.NewInt(0).Sub(recOngB2, recOngB1)
		log.Debugf("recipient: %s ong increment %+v\n", recipient.ToBase58(), ongIncrement)

	}

	return nil
}



func (this *TestEnv) tokenToOntOutput(ongBought uint64, maxTokens *big.Int, invoker *ontology_go_sdk.Account, recipient common.Address) error {
	if err := this.refreshAcctBalance(); err != nil {
		return fmt.Errorf("tokenToOngOutput, refreshBalance err: %v", err)
	}

	exOngBalance1 := this.OnChainEState[0].OntdLiquid
	exTokenB1 := this.OnChainEState[0].TokenLiquid

	recOngB1 := this.OntdBalance[recipient]
	// Condition check
	if this.OnChainEState[0].OntdLiquid.Cmp(big.NewInt(0).SetUint64(ongBought)) < 0 {
		return fmt.Errorf("tokenToOngOutput, exchange ong balance: %v < ongBought: %v", this.OnChainEState[0].OntdLiquid, ongBought)
	}
	if this.OnChainTState[0].Balances[invoker.Address].Cmp(maxTokens) < 0 {
		return fmt.Errorf("tokenToOngOutput, invoker: %s, not have enough token balance", invoker.Address.ToBase58())
	}
	if this.OnChainTState[0].Allowances[invoker.Address].Cmp(maxTokens) < 0 {
		log.Debugf("tokenToOngOutput, invoker: %s, not have allowance token", invoker.Address.ToBase58())
		// approve token to exchange from invoker
		approveTxHash, err := this.Sdk.NeoVM.InvokeNeoVMContract(this.GasPrice, this.GasLimit, invoker, invoker, this.OnChainTState[0].TokenAddr, []interface{}{"approve", []interface{}{
			invoker.Address, this.OnChainEState[0].ExchangeAddr, maxTokens,
		}})
		if err != nil {
			return fmt.Errorf("tokenToOngOutput: %s, approve token to exchange err: %v", invoker.Address.ToBase58(), err)
		}
		if _, err := this.Sdk.WaitForGenerateBlock(this.WaitTxTimeOut, 1); err != nil {
			return fmt.Errorf("Ontology, not generate block after %+v, err: %v", this.WaitTxTimeOut, err)
		}
		utils.PrintSmartEventByHash_Ont(this.Sdk, approveTxHash.ToHexString())
	}

	var params []interface{}
	// tokenToOngInput
	if invoker.Address == recipient {
		params = []interface{}{
			"tokenToOntSwapOutput",
			[]interface{}{
				ongBought,
				maxTokens,
				time.Now().Add(this.WaitTxTimeOut).Unix(),
				invoker.Address,
			},
		}
	} else {
		params = []interface{}{
			"tokenToOntTransferOutput",
			[]interface{}{
				ongBought,
				maxTokens,
				time.Now().Add(this.WaitTxTimeOut).Unix(),
				recipient,
				invoker.Address,
			},
		}
	}
	txHash, err := this.Sdk.NeoVM.InvokeNeoVMContract(this.GasPrice, this.GasLimit, invoker, invoker, this.OnChainEState[0].ExchangeAddr, params)
	if err != nil {
		return fmt.Errorf("tokenToOngOutput, invoker: %s invoke err: %v", invoker.Address.ToBase58(), err)
	}
	if _, err := this.Sdk.WaitForGenerateBlock(this.WaitTxTimeOut, 1); err != nil {
		return fmt.Errorf("Ontology, not generate block after %+v, err: %v", this.WaitTxTimeOut, err)
	}
	utils.PrintSmartEventByHash_Ont(this.Sdk, txHash.ToHexString())



	if err := this.refreshAcctBalance(); err != nil {
		return fmt.Errorf("tokenToOngInput, refreshBalance err: %v", err)
	}
	exOngBalance2 := this.OnChainEState[0].OntdLiquid
	exTokenB2 := this.OnChainEState[0].TokenLiquid
	exOngInc := big.NewInt(0).Sub(exOngBalance2, exOngBalance1)
	exTokenInc := big.NewInt(0).Sub(exTokenB2, exTokenB1)
	log.Debugf("exchange ong increased by : %v\n", exOngInc)
	log.Debugf("exchange token increased by : %v\n", exTokenInc)

	if invoker.Address == recipient {

	} else {
		recOngB2 := this.OntdBalance[recipient]
		ongIncrement := big.NewInt(0).Sub(recOngB2, recOngB1)
		log.Debugf("recipient: %s received %+v ong\n", recipient.ToBase58(), ongIncrement)
	}

	return nil
}


func (this *TestEnv) tokenToTokenInput(tokenSoldIndex int, tokenSold *big.Int, minTokenBought *big.Int, minOntdBought *big.Int, invoker *ontology_go_sdk.Account, recipient, tokenAddr common.Address) error {
	if err := this.refreshAcctBalance(); err != nil {
		return fmt.Errorf("tokenToTokenInput, refreshBalance err: %v", err)
	}

	exOngBalance1 := this.OnChainEState[0].OntdLiquid
	exTokenB1 := this.OnChainEState[0].TokenLiquid
	recOngB1 := this.OnChainTState[0].Balances[recipient]

	// Condition check
	if this.OnChainEState[0].TokenLiquid.Cmp(tokenSold) < 0 {
		return fmt.Errorf("tokenToTokenInput, exchange token balance: %v < tokenSold: %v", this.OnChainEState[0].TokenLiquid.String(), tokenSold.String())
	}
	if this.OnChainTState[0].Balances[invoker.Address].Cmp(tokenSold) < 1 {
		return fmt.Errorf("tokenToTokenInput, invoker: %s, not have enough token balance", invoker.Address.ToBase58())
	}
	if this.OnChainTState[0].Allowances[invoker.Address].Cmp(tokenSold) < 1 {
		log.Debugf("tokenToTokenInput, invoker: %s, not have allowance token", invoker.Address.ToBase58())
		// approve token to exchange from invoker
		approveTxHash, err := this.Sdk.NeoVM.InvokeNeoVMContract(this.GasPrice, this.GasLimit, invoker, invoker, this.OnChainTState[0].TokenAddr, []interface{}{"approve", []interface{}{
			invoker.Address, this.OnChainEState[0].ExchangeAddr, tokenSold,
		}})
		if err != nil {
			return fmt.Errorf("tokenToTokenInput: %s, approve token to exchange err: %v", invoker.Address.ToBase58(), err)
		}
		if _, err := this.Sdk.WaitForGenerateBlock(this.WaitTxTimeOut, 1); err != nil {
			return fmt.Errorf("Ontology, not generate block after %+v, err: %v", this.WaitTxTimeOut, err)
		}
		utils.PrintSmartEventByHash_Ont(this.Sdk, approveTxHash.ToHexString())
	}

	var params []interface{}
	// tokenToOngInput
	if invoker.Address == recipient {
		params = []interface{}{
			"tokenToTokenSwapInput",
			[]interface{}{
				tokenSold,
				minTokenBought,
				minOntdBought,
				time.Now().Add(this.WaitTxTimeOut).Unix(),
				tokenAddr,
				invoker.Address,
			},
		}
	} else {
		params = []interface{}{
			"tokenToTokenTransferInput",
			[]interface{}{
				tokenSold,
				minTokenBought,
				minOntdBought,
				time.Now().Add(this.WaitTxTimeOut).Unix(),
				recipient,
				tokenAddr,
				invoker.Address,
			},
		}
	}
	txHash, err := this.Sdk.NeoVM.InvokeNeoVMContract(this.GasPrice, this.GasLimit, invoker, invoker, this.OnChainEState[0].ExchangeAddr, params)
	if err != nil {
		return fmt.Errorf("tokenToTokenInput, invoker: %s invoke err: %v", invoker.Address.ToBase58(), err)
	}
	if _, err := this.Sdk.WaitForGenerateBlock(this.WaitTxTimeOut, 1); err != nil {
		return fmt.Errorf("Ontology, not generate block after %+v, err: %v", this.WaitTxTimeOut, err)
	}
	utils.PrintSmartEventByHash_Ont(this.Sdk, txHash.ToHexString())



	if err := this.refreshAcctBalance(); err != nil {
		return fmt.Errorf("tokenToOngInput, refreshBalance err: %v", err)
	}
	exOngBalance2 := this.OnChainEState[0].OntdLiquid
	exTokenB2 := this.OnChainEState[0].TokenLiquid
	exOngInc := big.NewInt(0).Sub(exOngBalance2, exOngBalance1)
	exTokenInc := big.NewInt(0).Sub(exTokenB2, exTokenB1)
	log.Debugf("exchange ong increased by : %s\n", exOngInc.String())
	log.Debugf("exchange token increased by : %s\n", exTokenInc.String())

	if invoker.Address == recipient {

	} else {
		recOngB2 := this.OnChainTState[0].Balances[recipient]
		ongIncrement := big.NewInt(0).Sub(recOngB2, recOngB1)
		log.Debugf("recipient: %s ong increment %+v\n", recipient.ToBase58(), ongIncrement)

	}
	return nil
}



func (this *TestEnv) tokenToTokenOutput(tokenBoughtIndex int, tokenBought *big.Int, maxTokenSold *big.Int, maxOntdSold *big.Int, invoker *ontology_go_sdk.Account, recipient, tokenAddr common.Address) error {
	if err := this.refreshAcctBalance(); err != nil {
		return fmt.Errorf("tokenToTokenInput, refreshBalance err: %v", err)
	}

	exOngBalance1 := this.OnChainEState[0].OntdLiquid
	exTokenB1 := this.OnChainEState[0].TokenLiquid
	recOngB1 := this.OnChainTState[0].Balances[recipient]

	// Condition check
	if this.OnChainEState[0].TokenLiquid.Cmp(tokenBought) < 0 {
		return fmt.Errorf("tokenToTokenOutput, exchange token balance: %v < tokenSold: %v", this.OnChainEState[0].TokenLiquid.String(), tokenBought.String())
	}
	if this.OnChainTState[0].Balances[invoker.Address].Cmp(maxTokenSold) < 1 {
		return fmt.Errorf("tokenToTokenOutput, invoker: %s, not have enough token balance", invoker.Address.ToBase58())
	}
	if this.OnChainTState[0].Allowances[invoker.Address].Cmp(maxTokenSold) < 1 {
		log.Debugf("tokenToTokenOutput, invoker: %s, not have allowance token", invoker.Address.ToBase58())
		// approve token to exchange from invoker
		approveTxHash, err := this.Sdk.NeoVM.InvokeNeoVMContract(this.GasPrice, this.GasLimit, invoker, invoker, this.OnChainTState[0].TokenAddr, []interface{}{"approve", []interface{}{
			invoker.Address, this.OnChainEState[0].ExchangeAddr, maxTokenSold,
		}})
		if err != nil {
			return fmt.Errorf("tokenToTokenOutput: %s, approve token to exchange err: %v", invoker.Address.ToBase58(), err)
		}
		if _, err := this.Sdk.WaitForGenerateBlock(this.WaitTxTimeOut, 1); err != nil {
			return fmt.Errorf("Ontology, not generate block after %+v, err: %v", this.WaitTxTimeOut, err)
		}
		utils.PrintSmartEventByHash_Ont(this.Sdk, approveTxHash.ToHexString())
	}

	var params []interface{}
	// tokenToOngInput
	if invoker.Address == recipient {
		params = []interface{}{
			"tokenToTokenSwapOutput",
			[]interface{}{
				tokenBought,
				maxTokenSold,
				maxOntdSold,
				time.Now().Add(this.WaitTxTimeOut).Unix(),
				tokenAddr,
				invoker.Address,
			},
		}
	} else {
		params = []interface{}{
			"tokenToTokenTransferOutput",
			[]interface{}{
				tokenBought,
				maxTokenSold,
				maxOntdSold,
				time.Now().Add(this.WaitTxTimeOut).Unix(),
				recipient,
				tokenAddr,
				invoker.Address,
			},
		}
	}
	txHash, err := this.Sdk.NeoVM.InvokeNeoVMContract(this.GasPrice, this.GasLimit, invoker, invoker, this.OnChainEState[0].ExchangeAddr, params)
	if err != nil {
		return fmt.Errorf("tokenToTokenOutput, invoker: %s invoke err: %v", invoker.Address.ToBase58(), err)
	}
	if _, err := this.Sdk.WaitForGenerateBlock(this.WaitTxTimeOut, 1); err != nil {
		return fmt.Errorf("Ontology, not generate block after %+v, err: %v", this.WaitTxTimeOut, err)
	}
	utils.PrintSmartEventByHash_Ont(this.Sdk, txHash.ToHexString())



	if err := this.refreshAcctBalance(); err != nil {
		return fmt.Errorf("tokenToTokenOutput, refreshBalance err: %v", err)
	}
	exOngBalance2 := this.OnChainEState[0].OntdLiquid
	exTokenB2 := this.OnChainEState[0].TokenLiquid
	exOngInc := big.NewInt(0).Sub(exOngBalance2, exOngBalance1)
	exTokenInc := big.NewInt(0).Sub(exTokenB2, exTokenB1)
	log.Debugf("exchange ong increased by : %s\n", exOngInc.String())
	log.Debugf("exchange token increased by : %s\n", exTokenInc.String())

	if invoker.Address == recipient {

	} else {
		recOngB2 := this.OnChainTState[0].Balances[recipient]
		ongIncrement := big.NewInt(0).Sub(recOngB2, recOngB1)
		log.Debugf("recipient: %s ong increment %+v\n", recipient.ToBase58(), ongIncrement)

	}
	return nil
}



func (this *TestEnv) tokenToExchangeInput(tokenSoldIndex int, tokenSold *big.Int, minTokenBought *big.Int, minOntdBought *big.Int, invoker *ontology_go_sdk.Account, recipient, exAddr common.Address) error {
	if err := this.refreshAcctBalance(); err != nil {
		return fmt.Errorf("tokenToExchangeInput, refreshBalance err: %v", err)
	}

	exOngBalance1 := this.OnChainEState[0].OntdLiquid
	exTokenB1 := this.OnChainEState[0].TokenLiquid
	recOngB1 := this.OnChainTState[0].Balances[recipient]

	// Condition check
	if this.OnChainEState[0].TokenLiquid.Cmp(tokenSold) < 0 {
		return fmt.Errorf("tokenToExchangeInput, exchange token balance: %v < tokenSold: %v", this.OnChainEState[0].TokenLiquid.String(), tokenSold.String())
	}
	if this.OnChainTState[0].Balances[invoker.Address].Cmp(tokenSold) < 1 {
		return fmt.Errorf("tokenToExchangeInput, invoker: %s, not have enough token balance", invoker.Address.ToBase58())
	}
	if this.OnChainTState[0].Allowances[invoker.Address].Cmp(tokenSold) < 1 {
		log.Debugf("tokenToExchangeInput, invoker: %s, not have allowance token", invoker.Address.ToBase58())
		// approve token to exchange from invoker
		approveTxHash, err := this.Sdk.NeoVM.InvokeNeoVMContract(this.GasPrice, this.GasLimit, invoker, invoker, this.OnChainTState[0].TokenAddr, []interface{}{"approve", []interface{}{
			invoker.Address, this.OnChainEState[0].ExchangeAddr, tokenSold,
		}})
		if err != nil {
			return fmt.Errorf("tokenToExchangeInput: %s, approve token to exchange err: %v", invoker.Address.ToBase58(), err)
		}
		if _, err := this.Sdk.WaitForGenerateBlock(this.WaitTxTimeOut, 1); err != nil {
			return fmt.Errorf("Ontology, not generate block after %+v, err: %v", this.WaitTxTimeOut, err)
		}
		utils.PrintSmartEventByHash_Ont(this.Sdk, approveTxHash.ToHexString())
	}

	var params []interface{}
	// tokenToOngInput
	if invoker.Address == recipient {
		params = []interface{}{
			"tokenToExchangeSwapInput",
			[]interface{}{
				tokenSold,
				minTokenBought,
				minOntdBought,
				time.Now().Add(this.WaitTxTimeOut).Unix(),
				exAddr,
				invoker.Address,
			},
		}
	} else {
		params = []interface{}{
			"tokenToExchangeTransferInput",
			[]interface{}{
				tokenSold,
				minTokenBought,
				minOntdBought,
				time.Now().Add(this.WaitTxTimeOut).Unix(),
				recipient,
				exAddr,
				invoker.Address,
			},
		}
	}
	txHash, err := this.Sdk.NeoVM.InvokeNeoVMContract(this.GasPrice, this.GasLimit, invoker, invoker, this.OnChainEState[0].ExchangeAddr, params)
	if err != nil {
		return fmt.Errorf("tokenToExchangeInput, invoker: %s invoke err: %v", invoker.Address.ToBase58(), err)
	}
	if _, err := this.Sdk.WaitForGenerateBlock(this.WaitTxTimeOut, 1); err != nil {
		return fmt.Errorf("Ontology, not generate block after %+v, err: %v", this.WaitTxTimeOut, err)
	}
	utils.PrintSmartEventByHash_Ont(this.Sdk, txHash.ToHexString())



	if err := this.refreshAcctBalance(); err != nil {
		return fmt.Errorf("tokenToExchangeInput, refreshBalance err: %v", err)
	}
	exOngBalance2 := this.OnChainEState[0].OntdLiquid
	exTokenB2 := this.OnChainEState[0].TokenLiquid
	exOngInc := big.NewInt(0).Sub(exOngBalance2, exOngBalance1)
	exTokenInc := big.NewInt(0).Sub(exTokenB2, exTokenB1)
	log.Debugf("exchange ong increased by : %s\n", exOngInc.String())
	log.Debugf("exchange token increased by : %s\n", exTokenInc.String())

	if invoker.Address == recipient {

	} else {
		recOngB2 := this.OnChainTState[0].Balances[recipient]
		ongIncrement := big.NewInt(0).Sub(recOngB2, recOngB1)
		log.Debugf("recipient: %s ong increment %+v\n", recipient.ToBase58(), ongIncrement)

	}
	return nil
}


func (this *TestEnv) tokenToExchangeOutput(tokenBoughtIndex int, tokenBought *big.Int, maxTokenSold *big.Int, maxOntdSold *big.Int, invoker *ontology_go_sdk.Account, recipient, tokenAddr common.Address) error {
	if err := this.refreshAcctBalance(); err != nil {
		return fmt.Errorf("tokenToTokenInput, refreshBalance err: %v", err)
	}

	exOngBalance1 := this.OnChainEState[0].OntdLiquid
	exTokenB1 := this.OnChainEState[0].TokenLiquid
	recOngB1 := this.OnChainTState[0].Balances[recipient]

	// Condition check
	if this.OnChainEState[0].TokenLiquid.Cmp(tokenBought) < 0 {
		return fmt.Errorf("tokenToTokenOutput, exchange token balance: %v < tokenSold: %v", this.OnChainEState[0].TokenLiquid.String(), tokenBought.String())
	}
	if this.OnChainTState[0].Balances[invoker.Address].Cmp(maxTokenSold) < 1 {
		return fmt.Errorf("tokenToTokenOutput, invoker: %s, not have enough token balance", invoker.Address.ToBase58())
	}
	if this.OnChainTState[0].Allowances[invoker.Address].Cmp(maxTokenSold) < 1 {
		log.Debugf("tokenToTokenOutput, invoker: %s, not have allowance token", invoker.Address.ToBase58())
		// approve token to exchange from invoker
		approveTxHash, err := this.Sdk.NeoVM.InvokeNeoVMContract(this.GasPrice, this.GasLimit, invoker, invoker, this.OnChainTState[0].TokenAddr, []interface{}{"approve", []interface{}{
			invoker.Address, this.OnChainEState[0].ExchangeAddr, maxTokenSold,
		}})
		if err != nil {
			return fmt.Errorf("tokenToTokenOutput: %s, approve token to exchange err: %v", invoker.Address.ToBase58(), err)
		}
		if _, err := this.Sdk.WaitForGenerateBlock(this.WaitTxTimeOut, 1); err != nil {
			return fmt.Errorf("Ontology, not generate block after %+v, err: %v", this.WaitTxTimeOut, err)
		}
		utils.PrintSmartEventByHash_Ont(this.Sdk, approveTxHash.ToHexString())
	}

	var params []interface{}
	// tokenToOngInput
	if invoker.Address == recipient {
		params = []interface{}{
			"tokenToExchangeSwapOutput",
			[]interface{}{
				tokenBought,
				maxTokenSold,
				maxOntdSold,
				time.Now().Add(this.WaitTxTimeOut).Unix(),
				tokenAddr,
				invoker.Address,
			},
		}
	} else {
		params = []interface{}{
			"tokenToExchangeTransferOutput",
			[]interface{}{
				tokenBought,
				maxTokenSold,
				maxOntdSold,
				time.Now().Add(this.WaitTxTimeOut).Unix(),
				recipient,
				tokenAddr,
				invoker.Address,
			},
		}
	}
	txHash, err := this.Sdk.NeoVM.InvokeNeoVMContract(this.GasPrice, this.GasLimit, invoker, invoker, this.OnChainEState[0].ExchangeAddr, params)
	if err != nil {
		return fmt.Errorf("tokenToTokenOutput, invoker: %s invoke err: %v", invoker.Address.ToBase58(), err)
	}
	if _, err := this.Sdk.WaitForGenerateBlock(this.WaitTxTimeOut, 1); err != nil {
		return fmt.Errorf("Ontology, not generate block after %+v, err: %v", this.WaitTxTimeOut, err)
	}
	utils.PrintSmartEventByHash_Ont(this.Sdk, txHash.ToHexString())



	if err := this.refreshAcctBalance(); err != nil {
		return fmt.Errorf("tokenToTokenOutput, refreshBalance err: %v", err)
	}
	exOngBalance2 := this.OnChainEState[0].OntdLiquid
	exTokenB2 := this.OnChainEState[0].TokenLiquid
	exOngInc := big.NewInt(0).Sub(exOngBalance2, exOngBalance1)
	exTokenInc := big.NewInt(0).Sub(exTokenB2, exTokenB1)
	log.Debugf("exchange ong increased by : %s\n", exOngInc.String())
	log.Debugf("exchange token increased by : %s\n", exTokenInc.String())

	if invoker.Address == recipient {

	} else {
		recOngB2 := this.OnChainTState[0].Balances[recipient]
		ongIncrement := big.NewInt(0).Sub(recOngB2, recOngB1)
		log.Debugf("recipient: %s ong increment %+v\n", recipient.ToBase58(), ongIncrement)

	}
	return nil
}
