package exchange

import (
	"github.com/skyinglyh1/uniswap_v1_test/log"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestInit(t *testing.T) {
	// all account get init token from account[4]
	for i := 0; i < len(testEnv.OnChainEState[0].Providers); i++ {
		TransferTokenOtherAccount(testEnv.OnChainEState[0].TokenAddr, 0, 4, i)
		TransferTokenOtherAccount(testEnv.OnChainEState[1].TokenAddr, 1, 4, i)
		//providerAddr := testEnv.OnChainEState[0].Providers[i].Address
		//log.Infof("account: %s, ontDBalance: %+v, tokenBalance: %+v, exchange 1 shareBalance: %+v\n", providerAddr.ToBase58(), testEnv.OntdBalance[providerAddr], testEnv.OnChainTState[0].Balances[providerAddr], testEnv.OnChainEState[0].ShareBalance[providerAddr])
	}
}

func GetAccountBalance() map[string][]*big.Int {
	var balances = make(map[string][]*big.Int)
	for i := 0; i < len(testEnv.OnChainEState[0].Providers); i++ {
		providerAddr := testEnv.OnChainEState[0].Providers[i].Address
		balances[providerAddr.ToBase58()] = append(balances[providerAddr.ToBase58()], testEnv.OntdBalance[providerAddr], testEnv.OnChainTState[0].Balances[providerAddr], testEnv.OnChainEState[0].ShareBalance[providerAddr])
		log.Infof("account: %s, ontDBalance: %+v, tokenBalance: %+v, exchange 1 shareBalance: %+v\n", providerAddr.ToBase58(), testEnv.OntdBalance[providerAddr], testEnv.OnChainTState[0].Balances[providerAddr], testEnv.OnChainEState[0].ShareBalance[providerAddr])
	}
	return balances
}

// just one address add and remove one by one
func TestRemoveAllLiquidity(t *testing.T) {
	//providerAddr := testEnv.OnChainEState[0].Providers[0].Address
	if err := testEnv.addLiquid(0, big.NewInt(100), big.NewInt(0).Add(big.NewInt(100000), big.NewInt(100000)), big.NewInt(200000)); err != nil {
		log.Errorf("address: %s, addLiquid() error: %+v", err)
	}
	for i := 0; i < len(testEnv.OnChainEState[0].Providers); i++ {
		amount1 := getShare(0, i)
		log.Infof("shareBalanceRes is %x ", amount1)
		if err := testEnv.removeLiquid(0, amount1, big.NewInt(1), testEnv.OnChainEState[0].Providers[i]); err != nil {
			log.Errorf("address: %s, removeLiquid() error: %+v", err)
		}
		amount2 := getShare(0, i)
		log.Infof("get shareBalance error %d ", amount2)
		assert.Equal(t, big.NewInt(0), amount2)
	}
	//  TODO: check exchange lock ontd is 0

}

// 4 account addLiquidity and one by one remove all
func TestRemoveAllLiquidity1In4(t *testing.T) {
	//all account add
	if err := testEnv.addLiquid(0, big.NewInt(100), big.NewInt(0).Add(big.NewInt(100000), big.NewInt(100000)), big.NewInt(200000)); err != nil {
		log.Errorf("address: %s, addLiquid() error: %+v", err)
	}
	if err := testEnv.addLiquid(1, big.NewInt(100), big.NewInt(0).Add(big.NewInt(100000), big.NewInt(100000)), big.NewInt(200000)); err != nil {
		log.Errorf("address: %s, addLiquid() error: %+v", err)
	}
	amount1 := getShare(0, 0)
	log.Infof("shareBalanceRes is %x ", amount1)
	if err := testEnv.removeLiquid(0, amount1, big.NewInt(1), testEnv.OnChainEState[0].Providers[0]); err != nil {
		log.Errorf("address: %s, removeLiquid() error: %+v", err)
	}
	amount2 := getShare(0, 0)
	log.Infof("get shareBalance  %d ", amount2)

	assert.Equal(t, big.NewInt(0), amount2)
}

// remove 50%
func TestRemoveLiquidity(t *testing.T) {
	//all account add
	if err := testEnv.addLiquid(0, big.NewInt(100), big.NewInt(0).Add(big.NewInt(100000), big.NewInt(100000)), big.NewInt(200000)); err != nil {
		log.Errorf("address: %s, addLiquid() error: %+v", err)
	}
	amount1 := getShare(0, 0)
	log.Infof("shareBalanceRes is %x ", amount1)
	removeAount := amount1.Div(amount1, big.NewInt(2))
	if err := testEnv.removeLiquid(0, removeAount, big.NewInt(1), testEnv.OnChainEState[0].Providers[0]); err != nil {
		log.Errorf("address: %s, removeLiquid() error: %+v", err)
	}
	amount2 := getShare(0, 0)
	log.Infof("get shareBalance  %d ", amount2)
	assert.Equal(t, removeAount, amount2)
}

func getShare(exhcangeIndex int, accountIndex int) *big.Int {
	providerAddr := testEnv.OnChainEState[exhcangeIndex].Providers[accountIndex].Address
	shareBalanceRes, err := testEnv.Sdk.NeoVM.PreExecInvokeNeoVMContract(testEnv.OffChainEState[exhcangeIndex].ExchangeAddr, []interface{}{"balanceOf", []interface{}{providerAddr}})
	if err != nil {
		log.Errorf("get shareBalance error %s ", err)
	}
	amount, _ := shareBalanceRes.Result.ToInteger()
	return amount
}

func TestTransferShare(t *testing.T) {
	// all account get init token from account[4]
	for i := 0; i < len(testEnv.OnChainEState[0].Providers); i++ {
		TransferTokenOtherAccount(testEnv.OnChainEState[0].TokenAddr, 0, 4, i)
		TransferTokenOtherAccount(testEnv.OnChainEState[1].TokenAddr, 1, 4, i)
		//providerAddr := testEnv.OnChainEState[0].Providers[i].Address
		//log.Infof("account: %s, ontDBalance: %+v, tokenBalance: %+v, exchange 1 shareBalance: %+v\n", providerAddr.ToBase58(), testEnv.OntdBalance[providerAddr], testEnv.OnChainTState[0].Balances[providerAddr], testEnv.OnChainEState[0].ShareBalance[providerAddr])
	}
}
