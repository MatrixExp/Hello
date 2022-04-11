package eth

import (
	"testing"
	"time"
)

const (
	ethRpcUrl = "https://data-seed-prebsc-1-s1.binance.org:8545"

	// scan https://etherscan.io/
	ethMainProdRpcUrl = "https://mainnet.infura.io/v3/da3717f25f824cc1baa32d812386d93f"

	// scan https://testnet.bscscan.com/
	binanceTestRpcUrl = "https://data-seed-prebsc-1-s1.binance.org:8545"

	// scan https://bscscan.com
	binanceProdRpcUrl = "https://bsc-dataseed.binance.org"

	// scan https://scan.sherpax.io/
	sherpaxProdRpcUrl = "https://mainnet.sherpax.io/rpc"
)

// var ethChain, _ = NewEthChain().CreateRemote(ethRpcUrl)

const (
	contractUSDT    = "0xdac17f958d2ee523a2206206994597c13d831ec7"
	contractBSCUSDT = "0x55d398326f99059fF775485246999027B3197955"
	contractBSCBUSD = "0xe9e7CEA3DedcA5984780Bafc599bD69ADd087D56"
	contractBSCUSDC = "0x8ac76a51cc950d9822d68b83fe1ad97b32cd580d"

	contractWKSXUSB  = "0xE7e312dfC08e060cda1AF38C234AEAcc7A982143"
	contractWKSXUSDT = "0x4B53739D798EF0BEa5607c254336b40a93c75b52"
	contractWKSXBUSD = "0x37088186089c7D6BcD556d9A15087DFaE3Ba0C32"
	contractWKSXUSDC = "0x935CC842f220CF3A7D10DA1c99F01B1A6894F7C5"
)

func TestConnect(t *testing.T) {
	// errRpc := binanceTestRpcUrl
	// chain, err := NewEthChain().CreateRemote(errRpc)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(chain)
	defer func() {
		err := recover()
		if err != nil {
			t.Log("recover err ", err)
		}
	}()

	chain, _ := NewEthChain().CreateRemote(binanceTestRpcUrl)

	for i := 0; i < 100; i++ {
		time.Sleep(1 * time.Second)
		address := "0xed24fc36d5ee211ea25a80239fb8c4cfd80f12ee"
		balance, err := chain.TokenBalance(address, address)
		if err != nil {
			t.Log("...... catched err", err)
		} else {
			t.Log("...... balance", balance)
		}
	}

	// t.Log("should successd connect", chain)
}
