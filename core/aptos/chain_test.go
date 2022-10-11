package aptos

import (
	"encoding/hex"
	"reflect"
	"testing"

	txbuilder "github.com/coming-chat/go-aptos/transaction_builder"
	"github.com/coming-chat/lcs"
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/core/testcase"
)

const (
	devnetRestUrl  = "https://fullnode.devnet.aptoslabs.com"
	testnetRestUrl = "https://testnet.aptoslabs.com"
)

func TestFaucet(t *testing.T) {
	account, _ := NewAccountWithMnemonic(testcase.M1)
	hashs, err := FaucetFundAccount(account.Address(), 10000000000, "")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hashs.Value)
}

func TestEstimatePayloadGasFeeBCS(t *testing.T) {
	account, err := NewAccountWithMnemonic(testcase.M1)
	if err != nil {
		t.Fatal(err)
	}
	contractAddress := "0xb6d5bb1291ae2739b5341e860b8f42cd7e579a0d90057dba3651bc4d1492c7eb"
	chain := NewChainWithRestUrl(devnetRestUrl)

	var createABI = "0106637265617465b6d5bb1291ae2739b5341e860b8f42cd7e579a0d90057dba3651bc4d1492c7eb0a7265645f7061636b657400000205636f756e74020d746f74616c5f62616c616e636502"
	abiBytes := make([][]byte, 0)
	abiStrs := []string{createABI}
	for _, s := range abiStrs {
		bs, err := hex.DecodeString(s)
		if err != nil {
			t.Fatal(err)
		}
		abiBytes = append(abiBytes, bs)
	}
	abi, err := txbuilder.NewTransactionBuilderABI(abiBytes)
	if err != nil {
		t.Fatal(err)
	}
	functionName := contractAddress + "::red_packet::create"
	payloadAbi, err := abi.BuildTransactionPayload(
		functionName,
		[]string{},
		[]any{
			uint64(5), uint64(10000),
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	bs, err := lcs.Marshal(payloadAbi)
	if err != nil {
		t.Fatal(err)
	}
	fee, err := chain.EstimatePayloadGasFeeBCS(account, bs)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fee)
}

func TestFetchTransactionDetail(t *testing.T) {
	chain := NewChainWithRestUrl(devnetRestUrl)
	hash := "0x252e1bee1b5b534d44f59ca4ce3bb9742a454ca85692b0ef71894f8568fb9ef3"
	detail, err := chain.FetchTransactionDetail(hash)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(detail)
}

func TestChain_FetchTransactionDetail(t *testing.T) {
	tests := []struct {
		name    string
		restUrl string
		hash    string
		want    *base.TransactionDetail
		wantErr bool
	}{
		{
			name:    "testnet normal",
			restUrl: devnetRestUrl,
			hash:    "0xfd496b3dccae000096d4bf4aef581863ce2600c8867be9c2fe5b82a3408441af",
			want: &base.TransactionDetail{
				HashString:      "0xfd496b3dccae000096d4bf4aef581863ce2600c8867be9c2fe5b82a3408441af",
				FromAddress:     "0xcf4ddd208bbbbefb3227cafa5c917fc6541d26b1869276ea80d99ee0505fc6f8",
				ToAddress:       "0xcdbe33da8d218e97a9bec6443ba4a1b1858494f29142976d357f4770c384e015",
				Amount:          "100",
				EstimateFees:    "4",
				Status:          base.TransactionStatusSuccess,
				FinishTimestamp: 1660556054,
			},
		},
		{
			name:    "testnet failed tx",
			restUrl: devnetRestUrl,
			hash:    "0x0a5720b278708820eb46c24af485858da8668e183a27ee57f3eed402cdda7436",
			want: &base.TransactionDetail{
				HashString:      "0x0a5720b278708820eb46c24af485858da8668e183a27ee57f3eed402cdda7436",
				FromAddress:     "0x56252ac5625573224bcaf89119b046f4a35a5c23bbf3d0f3eaa43311fbd2b2b3",
				ToAddress:       "0x903056ed3ddd9c7b9d5231ac96c8e6a218fe2a7cf26f17f04a96edb2cb832566",
				Amount:          "100",
				EstimateFees:    "4",
				Status:          base.TransactionStatusFailure,
				FinishTimestamp: 1660617399,
				FailureMessage:  "Move abort by ECOIN_STORE_NOT_PUBLISHED\n When an account hasn't registered `CoinStore` for `CoinType`.",
			},
		},
		{
			name:    "set address to hash",
			restUrl: devnetRestUrl,
			hash:    "0xcf4ddd208bbbbefb3227cafa5c917fc6541d26b1869276ea80d99ee0505fc6f8",
			wantErr: true,
		},
		{
			name:    "not transfer",
			restUrl: devnetRestUrl,
			hash:    "0xb7c06248b83bb7854d75f8d09a56ce4f5d7f799445fdb8781fccc536a01cd971",
			wantErr: true,
		},
		{
			name:    "mint tx",
			restUrl: devnetRestUrl,
			hash:    "0x6934afd26b2e371f69ed2095dab30961b4c5c4b40fca2351966cbcd6add96a69",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewChainWithRestUrl(tt.restUrl)
			got, err := c.FetchTransactionDetail(tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("Chain.FetchTransactionDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Chain.FetchTransactionDetail() = %v, want %v", got, tt.want)
			}
		})
	}
}
