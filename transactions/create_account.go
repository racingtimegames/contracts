package transactions

import (
	"context"
	"contracts/utils/lib"
	"encoding/hex"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/templates"
)

func CreateAddress(client *client.Client, seed string, sigAlgo crypto.SignatureAlgorithm, hashAlgo crypto.HashAlgorithm, auth lib.Account) (accountAddr string, privateKey string, err error) {

	_privateKey, err := crypto.GeneratePrivateKey(sigAlgo, []byte(seed))
	if err != nil {
		return
	}

	privateKey = hex.EncodeToString(_privateKey.Encode())

	myAcctKey := flow.NewAccountKey().
		SetPublicKey(_privateKey.PublicKey()).
		SetSigAlgo(_privateKey.Algorithm()).
		SetHashAlgo(hashAlgo).
		SetWeight(flow.AccountKeyWeightThreshold)

	serviceAcctAddr, serviceAcctKey, serviceSigner, err := lib.ServiceAccount(client, auth.SigAlgo, auth.HashAlgo, auth.KeyIndex, auth.Address, auth.PrivateKey)
	if err != nil {
		return
	}

	referenceBlockID := lib.GetReferenceBlockId(client)

	createAccountTx, err := templates.CreateAccount([]*flow.AccountKey{myAcctKey}, nil, serviceAcctAddr)
	createAccountTx.SetProposalKey(
		serviceAcctAddr,
		serviceAcctKey.Index,
		serviceAcctKey.SequenceNumber,
	)
	createAccountTx.SetReferenceBlockID(referenceBlockID)
	createAccountTx.SetPayer(serviceAcctAddr)

	if err = createAccountTx.SignEnvelope(serviceAcctAddr, serviceAcctKey.Index, serviceSigner); err != nil {
		return
	}

	if err = client.SendTransaction(context.Background(), *createAccountTx); err != nil {
		return
	}

	accountCreationTxRes, err := lib.WaitForSeal(context.Background(), client, createAccountTx.ID())
	if err != nil {
		return
	}
	if accountCreationTxRes.Error != nil {
		return "", "", accountCreationTxRes.Error
	}
	var _address flow.Address
	for _, event := range accountCreationTxRes.Events {
		if event.Type == flow.EventAccountCreated {
			accountCreatedEvent := flow.AccountCreatedEvent(event)
			_address = accountCreatedEvent.Address()
		}
	}
	accountAddr = "0x" + _address.Hex()
	return
}
