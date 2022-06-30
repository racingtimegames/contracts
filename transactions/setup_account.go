package transactions

import (
	"context"
	"contracts/utils/lib"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
)

func SetupAccount(client *client.Client, script string, auth lib.Account) (transactionId string, result *flow.TransactionResult, err error) {
	referenceBlock, err := client.GetLatestBlock(context.Background(), false)
	if err != nil {
		panic(err)
	}

	acctAddress, acctKey, signer, err := lib.ServiceAccount(client, auth.SigAlgo, auth.HashAlgo, auth.KeyIndex, auth.Address, auth.PrivateKey)
	if err != nil {
		return
	}

	tx := flow.NewTransaction().
		SetScript([]byte(script)).
		SetGasLimit(100).
		SetProposalKey(acctAddress, acctKey.Index, acctKey.SequenceNumber).
		SetReferenceBlockID(referenceBlock.ID).
		SetPayer(acctAddress).
		AddAuthorizer(acctAddress)

	if err = tx.SignEnvelope(acctAddress, acctKey.Index, signer); err != nil {
		return
	}

	if err = client.SendTransaction(context.Background(), *tx); err != nil {
		return
	}
	result, err = lib.WaitForSeal(context.Background(), client, tx.ID())
	if err != nil {
		return
	}

	return
}
