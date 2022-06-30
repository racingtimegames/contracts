package transactions

import (
	"context"
	"contracts/utils/lib"
	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
)

type NFTInfo struct {
	RewardID     uint32
	TypeID       uint32
	SerialNumber uint32
	Ipfs         string
	Address      string
}

func MintNFT(client *client.Client, script string, nft NFTInfo, auth lib.Account) (tx *flow.Transaction, result *flow.TransactionResult, err error) {
	//_script := fmt.Sprintf(mintArt, global.Config.Flow.NonFungibleToken, global.Config.Flow.RacingTime)

	referenceBlock, err := client.GetLatestBlock(context.Background(), false)
	if err != nil {
		return
	}

	acctAddress, acctKey, signer, err := lib.ServiceAccount(client, auth.SigAlgo, auth.HashAlgo, auth.KeyIndex, auth.Address, auth.PrivateKey)
	if err != nil {
		return
	}

	tx = flow.NewTransaction().
		SetScript([]byte(script)).
		SetGasLimit(100).
		SetProposalKey(acctAddress, acctKey.Index, acctKey.SequenceNumber).
		SetReferenceBlockID(referenceBlock.ID).
		SetPayer(acctAddress).
		AddAuthorizer(acctAddress)

	if err = tx.AddArgument(cadence.NewAddress(flow.HexToAddress(nft.Address))); err != nil {
		return
	}

	if err = tx.AddArgument(cadence.NewUInt32(nft.RewardID)); err != nil {
		return
	}

	if err = tx.AddArgument(cadence.NewUInt32(nft.TypeID)); err != nil {
		return
	}

	if err = tx.AddArgument(cadence.NewUInt32(nft.SerialNumber)); err != nil {
		return
	}
	_ipfs, err := cadence.NewString(nft.Ipfs)
	if err != nil {
		return
	}
	if err = tx.AddArgument(_ipfs); err != nil {
		return
	}

	if err = tx.SignEnvelope(acctAddress, acctKey.Index, signer); err != nil {
		return
	}

	if err = client.SendTransaction(context.Background(), *tx); err != nil {
		return
	}

	result, err = lib.WaitForSeal(context.Background(), client, tx.ID())
	return
}
