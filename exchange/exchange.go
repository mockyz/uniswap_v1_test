package exchange

import (
	ontology_go_sdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
	"github.com/skyinglyh1/uniswap_v1_test/config"
	"github.com/skyinglyh1/uniswap_v1_test/log"
	"github.com/skyinglyh1/uniswap_v1_test/utils"
	"os"
)

type ExchangeTest struct {
	Sdk *ontology_go_sdk.OntologySdk
	Accts []*ontology_go_sdk.Account
	FactoryHash common.Address
	ExchangeHash1 common.Address
	TokenHash1 common.Address
	ExchangeHash2 common.Address
	TokenHash2 common.Address
	TestMode uint64
}

func NewExchangeTest(config *config.Config) *ExchangeTest {
	sdk, accts, err := utils.GetSdkAndAccount(config.OntRpcAddress, config.WalletPath, config.AcctPwd)
	if err != nil {
		log.Errorf("GetSdkAndAccount err: %v", err)
		os.Exit(1)
	}
	factoryHash, err1 := common.AddressFromHexString(config.FactoryHash)
	ex1, err2 := common.AddressFromHexString(config.Exchange1Hash)
	t1, err3 := common.AddressFromHexString(config.Token1Hash)
	if err1 != nil || err2 != nil || err3 != nil {
		log.Errorf("FactoryHash err: %v, Exchange1Hash1err: %v, Token1Hash1 err: %v", err1, err2, err3)
		os.Exit(1)
	}
	et := &ExchangeTest{
		Sdk: sdk,
		Accts: accts,
		FactoryHash: factoryHash,
		ExchangeHash1: ex1,
		TokenHash1: t1,
		TestMode: config.TestFlag,
	}
	if config.TestFlag == 0 {

	} else {
		ex2, err1 := common.AddressFromHexString(config.Exchange2Hash)
		t2, err2 := common.AddressFromHexString(config.Token2Hash)
		if err1 != nil || err2 != nil{
			log.Errorf("Exchange1Hash2 err: %v, Token1Hash2 err: %v", err1, err2)
			os.Exit(1)
		}
		et.ExchangeHash2 = ex2
		et.TokenHash2 = t2
	}
	return et
}



func (this *ExchangeTest) Liquid() {

}

func (this *ExchangeTest) Trade() {
	
	go this.Token1ToOng()
	go this.OngToToken1()

	if (this.TestMode == 1) {
		go this.Token1ToToken2()
		go this.Token2ToToken1()
	}
}

func (this *ExchangeTest) OngToToken1() {

}

func (this *ExchangeTest) Token1ToOng() {

}

func (this *ExchangeTest) Token1ToToken2() {

}

func (this *ExchangeTest) Token2ToToken1() {

}