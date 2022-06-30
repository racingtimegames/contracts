package transactions

import (
	"context"
	"contracts/utils/lib"
	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
)

func BatchTransferRacingTimeNFT(client *client.Client, script string, toAddress string, nftId []uint64, auth lib.Account) (transaction *flow.Transaction, err error) {
	referenceBlock, err := client.GetLatestBlock(context.Background(), false)
	if err != nil {
		return
	}

	acctAddress, acctKey, signer, err := lib.ServiceAccount(client, auth.SigAlgo, auth.HashAlgo, auth.KeyIndex, auth.Address, auth.PrivateKey)
	if err != nil {
		return
	}

	transaction = flow.NewTransaction().
		SetScript([]byte(script)).
		SetGasLimit(9999).
		SetProposalKey(acctAddress, acctKey.Index, acctKey.SequenceNumber).
		SetReferenceBlockID(referenceBlock.ID).
		SetPayer(acctAddress).
		AddAuthorizer(acctAddress)

	var values []cadence.Value
	for _, id := range nftId {
		values = append(values, cadence.NewUInt64(id))
	}

	if err = transaction.AddArgument(cadence.NewAddress(flow.HexToAddress(toAddress))); err != nil {
		return
	}

	if err = transaction.AddArgument(cadence.NewArray(values)); err != nil {
		return
	}

	if err = transaction.SignEnvelope(acctAddress, acctKey.Index, signer); err != nil {
		return
	}

	if err = client.SendTransaction(context.Background(), *transaction); err != nil {
		return
	}

	return
}
