package aptos

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/go-aptos/aptosclient"
	"github.com/coming-chat/go-aptos/aptostypes"
	"github.com/coming-chat/wallet-SDK/core/base"
)

const (
	MaxGasAmount = 1000
	GasPrice     = 1
)

type TransactionOption interface {
	Process(tx *aptostypes.Transaction) (*aptostypes.Transaction, error)
}

type IChain interface {
	base.Chain
	SubmitTransactionPayload(account base.Account, payload *aptostypes.Payload, option TransactionOption) (string, error)
	GetClient() (*aptosclient.RestClient, error)
}

type Chain struct {
	restClient *aptosclient.RestClient
	RestUrl    string
}

func NewChainWithRestUrl(restUrl string) *Chain {
	return &Chain{RestUrl: restUrl}
}

func (c *Chain) client() (*aptosclient.RestClient, error) {
	if c.restClient != nil {
		return c.restClient, nil
	}
	var err error
	c.restClient, err = aptosclient.Dial(context.Background(), c.RestUrl)
	return c.restClient, err
}

// MARK - Implement the protocol Chain

func (c *Chain) MainToken() base.Token {
	return &Token{chain: c}
}

func (c *Chain) BalanceOfAddress(address string) (b *base.Balance, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	client, err := c.client()
	if err != nil {
		return
	}

	balance, err := client.BalanceOf(address)
	if err != nil {
		return
	}

	return &base.Balance{
		Total:  balance.String(),
		Usable: balance.String(),
	}, nil
}

func (c *Chain) BalanceOfPublicKey(publicKey string) (*base.Balance, error) {
	address, err := EncodePublicKeyToAddress(publicKey)
	if err != nil {
		return nil, err
	}
	return c.BalanceOfAddress(address)
}
func (c *Chain) BalanceOfAccount(account base.Account) (*base.Balance, error) {
	return c.BalanceOfAddress(account.Address())
}

// Send the raw transaction on-chain
// @return the hex hash string
func (c *Chain) SendRawTransaction(signedTx string) (hash string, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	bytes, err := types.HexDecodeString(signedTx)
	if err != nil {
		return
	}
	var transaction = &aptostypes.Transaction{}
	err = json.Unmarshal(bytes, transaction)
	if err != nil {
		return
	}

	client, err := c.client()
	if err != nil {
		return
	}
	resultTx, err := client.SubmitTransaction(transaction)
	if err != nil {
		return
	}

	return resultTx.Hash, nil
}

// Fetch transaction details through transaction hash
func (c *Chain) FetchTransactionDetail(hash string) (detail *base.TransactionDetail, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	client, err := c.client()
	if err != nil {
		return
	}

	transaction, err := client.GetTransactionByHash(hash)
	if err != nil {
		return
	}
	return toBaseTransaction(transaction)
}

func (c *Chain) FetchTransactionStatus(hash string) base.TransactionStatus {
	client, err := c.client()
	if err != nil {
		return base.TransactionStatusFailure
	}
	transaction, err := client.GetTransactionByHash(hash)
	if err != nil {
		return base.TransactionStatusFailure
	}
	if transaction.Success {
		return base.TransactionStatusSuccess
	} else {
		return base.TransactionStatusFailure
	}
}

func (c *Chain) BatchFetchTransactionStatus(hashListString string) string {
	hashList := strings.Split(hashListString, ",")
	statuses, _ := base.MapListConcurrentStringToString(hashList, func(s string) (string, error) {
		return strconv.Itoa(c.FetchTransactionStatus(s)), nil
	})
	return strings.Join(statuses, ",")
}

// MARK - Implement the protocol IChain

func (c *Chain) GetClient() (*aptosclient.RestClient, error) {
	return c.client()
}

func (c *Chain) SubmitTransactionPayload(account base.Account, payload *aptostypes.Payload, option TransactionOption) (string, error) {
	transaction, err := c.createTransactionFromPayload(account, payload)
	if err != nil {
		return "", err
	}

	if nil != option {
		transaction, err = option.Process(transaction)
		if err != nil {
			return "", err
		}
	}

	transaction, err = c.signTransaction(account, transaction)
	if err != nil {
		return "", err
	}

	client, err := c.client()
	if err != nil {
		return "", err
	}
	resultTx, err := client.SubmitTransaction(transaction)
	if err != nil {
		return "", err
	}

	return resultTx.Hash, nil
}

func (c *Chain) signTransaction(account base.Account, transaction *aptostypes.Transaction) (*aptostypes.Transaction, error) {
	client, err := c.client()
	if err != nil {
		return nil, err
	}
	signingMessage, err := client.CreateTransactionSigningMessage(transaction)
	if err != nil {
		return nil, err
	}
	signatureData, _ := account.Sign(signingMessage, "")
	transaction.Signature = &aptostypes.Signature{
		Type:      "ed25519_signature",
		PublicKey: account.PublicKeyHex(),
		Signature: types.HexEncodeToString(signatureData),
	}
	return transaction, nil
}

func (c *Chain) createTransactionFromPayload(account base.Account, payload *aptostypes.Payload) (*aptostypes.Transaction, error) {
	client, err := c.client()
	if err != nil {
		return nil, err
	}

	fromAddress := account.Address()
	accountData, err := client.GetAccount(fromAddress)
	if err != nil {
		return nil, err
	}
	ledgerInfo, err := client.LedgerInfo()
	if err != nil {
		return nil, err
	}

	return &aptostypes.Transaction{
		Sender:                  fromAddress,
		SequenceNumber:          accountData.SequenceNumber,
		MaxGasAmount:            MaxGasAmount,
		GasUnitPrice:            GasPrice,
		Payload:                 payload,
		ExpirationTimestampSecs: ledgerInfo.LedgerTimestamp + 600, // timeout 10 mins
	}, nil
}

/**
 * This creates an account if it does not exist and mints the specified amount of
 * coins into that account
 * @param address Hex-encoded 16 bytes Aptos account address wich mints tokens
 * @param amount Amount of tokens to mint
 * @param faucetUrl default https://faucet.devnet.aptoslabs.com
 * @returns Hashes of submitted transactions, e.g. "hash1,has2,hash3,..."
 */
func FaucetFundAccount(address string, amount int64, faucetUrl string) (h *base.OptionalString, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	hashs, err := aptosclient.FaucetFundAccount(address, uint64(amount), faucetUrl)
	if err != nil {
		return
	}
	return &base.OptionalString{Value: strings.Join(hashs[:], ",")}, nil
}

func toBaseTransaction(transaction *aptostypes.Transaction) (*base.TransactionDetail, error) {
	if transaction.Type != aptostypes.TypeUserTransaction ||
		transaction.Payload.Type != aptostypes.EntryFunctionPayload {
		return nil, errors.New("Invalid transfer transaction.")
	}

	detail := &base.TransactionDetail{
		HashString:  transaction.Hash,
		FromAddress: transaction.Sender,
	}

	gasFee := transaction.GasUnitPrice * transaction.GasUsed
	detail.EstimateFees = strconv.FormatUint(gasFee, 10)

	args := transaction.Payload.Arguments
	if len(args) >= 2 {
		detail.ToAddress = args[0].(string)
		detail.Amount = args[1].(string)
	}

	if transaction.Success {
		detail.Status = base.TransactionStatusSuccess
	} else {
		detail.Status = base.TransactionStatusFailure
		detail.FailureMessage = transaction.VmStatus
	}

	timestamp := transaction.Timestamp / 1e6
	detail.FinishTimestamp = int64(timestamp)

	return detail, nil
}
