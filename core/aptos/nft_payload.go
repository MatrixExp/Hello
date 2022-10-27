package aptos

import (
	"errors"
	"fmt"

	aptosnft "github.com/coming-chat/go-aptos/nft"
	txbuilder "github.com/coming-chat/go-aptos/transaction_builder"
	"github.com/coming-chat/lcs"
	"github.com/coming-chat/wallet-SDK/core/base"
)

/*
* BCS Payload builder for NFT (include Coming CID)

	### Demo
	example for `CIDTokenTransferPayload`
	```
	var payload, err = builder.CIDTokenTransferPayload(1234, receiverAddress)
	var gasPrice, err = chain.EstimateGasPrice()
	var gasAmount, err := chain.EstimatePayloadGasFeeBCS(account, payload)
	print("estimate gas fee = %s", gasPrice * gasAmount)
	var hash, err = chain.SubmitTransactionPayloadBCS(account, payload)
	print("submited hash = %s", hash)
	```
*/
type NFTPayloadBCSBuilder struct {
	// The contract provider for the CID module. Default is `0xc73d3c0a171c871e1858414734c7776f0b0cfa567c2af7f0070d1436aab2306b`
	CIDContract string
}

func NewNFTPayloadBCSBuilder(contractAddress string) *NFTPayloadBCSBuilder {
	if contractAddress == "" {
		contractAddress = "0xc73d3c0a171c871e1858414734c7776f0b0cfa567c2af7f0070d1436aab2306b"
	}
	return &NFTPayloadBCSBuilder{
		CIDContract: contractAddress,
	}
}

// MARK - Payload build for Coming CID

func (b *NFTPayloadBCSBuilder) cidModuleId() (*txbuilder.ModuleId, error) {
	address, err := txbuilder.NewAccountAddressFromHex(b.CIDContract)
	if err != nil {
		return nil, err
	}
	return &txbuilder.ModuleId{
		Address: *address,
		Name:    "cid",
	}, nil
}

func (b *NFTPayloadBCSBuilder) CIDAllowDirectTransferPayload() ([]byte, error) {
	module, err := b.cidModuleId()
	if err != nil {
		return nil, err
	}
	payload := txbuilder.TransactionPayloadEntryFunction{
		ModuleName:   *module,
		FunctionName: "allow_direct_transfer",
	}
	return lcs.Marshal(txbuilder.TransactionPayload(payload))
}

func (b *NFTPayloadBCSBuilder) CIDTokenTransferPayload(cid int64, toReceiver string) ([]byte, error) {
	module, err := b.cidModuleId()
	if err != nil {
		return nil, err
	}
	receiver, err := txbuilder.NewAccountAddressFromHex(toReceiver)
	if err != nil {
		return nil, errors.New("Receiver address is invalid.")
	}
	cidString := fmt.Sprintf("%v.aptos", cid)
	cidBytes := txbuilder.BCSSerializeBasicValue(cidString)
	payload := txbuilder.TransactionPayloadEntryFunction{
		ModuleName:   *module,
		FunctionName: "cid_token_transfer",
		Args: [][]byte{
			cidBytes, receiver[:],
		},
	}
	return lcs.Marshal(txbuilder.TransactionPayload(payload))
}

func (b *NFTPayloadBCSBuilder) CIDRegister(cid uint64) ([]byte, error) {
	module, err := b.cidModuleId()
	if err != nil {
		return nil, err
	}
	payload := txbuilder.TransactionPayloadEntryFunction{
		ModuleName:   *module,
		FunctionName: "register",
		Args: [][]byte{
			txbuilder.BCSSerializeBasicValue(cid),
		},
	}
	return lcs.Marshal(txbuilder.TransactionPayload(payload))
}

// MARK - Payload build for normal token (0x3::token)

func (b *NFTPayloadBCSBuilder) OfferTokenTransactionNFT(receiver string, nft *base.NFT) ([]byte, error) {
	return b.tokenPayloadBuild(receiver, nft.ContractAddress, nft.Collection, nft.Name, "offer")
}

/*
* Build payload that offer token

  - @param receiver The token receiver
  - @param creator The token creator
  - @param collection The token's collection name
  - @param name The token's name
  - @return The offer token payload.
*/
func (b *NFTPayloadBCSBuilder) OfferTokenTransactionParams(receiver, creator, collection, name string) ([]byte, error) {
	return b.tokenPayloadBuild(receiver, creator, collection, name, "offer")
}

func (b *NFTPayloadBCSBuilder) CancelTokenOffer(receiver, creator, collection, name string) ([]byte, error) {
	return b.tokenPayloadBuild(receiver, creator, collection, name, "cancelOffer")
}

/*
* Build payload that claim token, the nft info will be obtaining through offer hash.

  - @param offerHash The submitted hash of the transaction that offer the token
  - @param chain The chain on which the transaction resides
  - @param receiver The token receiver will be check whether it matches the nft offer information.
  - @return The claim token payload
*/
func (b *NFTPayloadBCSBuilder) ClaimTokenFromHash(offerHash string, chain *Chain, receiver string) (res []byte, err error) {
	client, err := chain.client()
	if err != nil {
		return
	}
	offeredTxn, err := client.GetTransactionByHash(offerHash)
	if err != nil {
		return
	}
	if !offeredTxn.Success {
		if offeredTxn.VmStatus != "" {
			return nil, errors.New("Claim failed, the offer transaction failed: " + offeredTxn.VmStatus)
		} else {
			return nil, errors.New("Claim failed, the offer transaction may still be pending.")
		}
		return //lint:ignore
	}
	if offeredTxn.Payload.Function != "0x3::token_transfers::offer_script" {
		return nil, errors.New("Claim failed, the given hash is not an offer token transaction")
	}
	arguments := offeredTxn.Payload.Arguments
	if len(arguments) < 4 {
		return nil, errors.New("Claim failed, offer params invalid.")
	}
	nftReceiver := arguments[0].(string)
	if receiver != nftReceiver {
		return nil, errors.New("Claim failed, this token is not offer to the receiver.")
	}
	creator := arguments[1].(string)
	collectionName := arguments[2].(string)
	tokenName := arguments[3].(string)
	nftSender := offeredTxn.Sender

	return b.tokenPayloadBuild(nftSender, creator, collectionName, tokenName, "claim")
}

/*
* Build payload that claim token

  - @param sender The transferred token owner
  - @param creator The token creator
  - @param collection The token's collection name
  - @param name The token's name
  - @return The claim token payload.
*/
func (b *NFTPayloadBCSBuilder) ClaimTokenTransactionParams(sender, creator, collection, name string) ([]byte, error) {
	return b.tokenPayloadBuild(sender, creator, collection, name, "claim")
}

// @param action enum of `offer`, `claim`, `cancelOffer`
func (b *NFTPayloadBCSBuilder) tokenPayloadBuild(senderOrReceiver, creator, collection, name string, action string) (res []byte, err error) {
	senderOrReceiverAddress, err := txbuilder.NewAccountAddressFromHex(senderOrReceiver)
	if err != nil {
		return
	}
	creatorAddress, err := txbuilder.NewAccountAddressFromHex(creator)
	if err != nil {
		return
	}
	if collection == "" || name == "" {
		return nil, errors.New("The `collection` and `name` cannot be empty.")
	}
	builder, err := aptosnft.NewNFTPayloadBuilder()
	if err != nil {
		return
	}

	var payload txbuilder.TransactionPayload
	switch action {
	case "offer":
		payload, err = builder.OfferToken(*senderOrReceiverAddress, *creatorAddress, collection, name, 1, 0)
	case "claim":
		payload, err = builder.ClaimToken(*senderOrReceiverAddress, *creatorAddress, collection, name, 0)
	case "cancelOffer":
		payload, err = builder.CancelTokenOffer(*senderOrReceiverAddress, *creatorAddress, collection, name, 0)
	default:
		return nil, errors.New("Invalid token action: " + action)
	}
	if err != nil {
		return
	}
	return lcs.Marshal(payload)
}

func (c *Chain) IsAllowedDirectTransferToken(account string) (*base.OptionalBool, error) {
	client, err := c.client()
	if err != nil {
		return nil, err
	}
	tokenStore, err := client.GetAccountResourceHandle404(account, "0x3::token::TokenStore", 0)
	if err != nil {
		return nil, err
	}
	if tokenStore == nil {
		return &base.OptionalBool{Value: false}, nil
	}
	if allow, ok := tokenStore.Data["direct_transfer"].(bool); ok {
		return &base.OptionalBool{Value: allow}, nil
	} else {
		return nil, errors.New("The queried data is incorrect.")
	}
}
