package lib

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"github.com/onflow/flow-go-sdk/crypto"
	"strconv"
	"time"
)

var (
	FAccount Account
)

type PrivateKey struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
}

type Account struct {
	SigAlgo    crypto.SignatureAlgorithm
	HashAlgo   crypto.HashAlgorithm
	KeyIndex   int
	Address    string
	PrivateKey string
}

func WaitForSeal(ctx context.Context, c *client.Client, id flow.Identifier) (result *flow.TransactionResult, err error) {
	result, _ = c.GetTransactionResult(ctx, id)
	for result.Status != flow.TransactionStatusSealed {
		time.Sleep(time.Second)
		result, err = c.GetTransactionResult(ctx, id)
		if err != nil {
			return
		}
		if result.Error != nil {
			return result, result.Error
		}
	}
	return
}

func ServiceAccount(flowClient *client.Client, sigAlgo crypto.SignatureAlgorithm, hashAlgo crypto.HashAlgorithm, keyIndex int, address, privateKey string) (flow.Address, *flow.AccountKey, crypto.Signer, error) {

	servicePrivateKeyHex := privateKey
	_privateKey, err := crypto.DecodePrivateKeyHex(sigAlgo, servicePrivateKeyHex)
	if err != nil {
		return [8]byte{}, nil, nil, err
	}

	_address := flow.HexToAddress(address)

	_account, err := flowClient.GetAccount(context.Background(), _address)
	if err != nil {
		return [8]byte{}, nil, nil, err
	}

	_accKey := _account.Keys[keyIndex]
	_signer, _err := crypto.NewInMemorySigner(_privateKey, hashAlgo)
	return _address, _accKey, _signer, _err
}

func GetAccountKey(flowClient *client.Client, public, address string) (key *flow.AccountKey, err error) {
	_address := flow.HexToAddress(address)

	_account, err := flowClient.GetAccount(context.Background(), _address)
	if err != nil {
		return nil, err
	}
	for _, key := range _account.Keys {
		if cadence.String(hex.EncodeToString(key.PublicKey.Encode())).String() == public {
			return key, nil
		}
	}
	return nil, errors.New("not find key")
}

func GetAccountKeyFromIndex(flowClient *client.Client, address string, index int) (*flow.AccountKey, error) {
	_address := flow.HexToAddress(address)

	_account, err := flowClient.GetAccount(context.Background(), _address)
	if err != nil {
		return nil, err
	}
	_key := _account.Keys[index]
	return _key, nil
}

func GetPublicKey(sigAlgo crypto.SignatureAlgorithm, privateKey string) crypto.PublicKey {
	_sigAlgo := crypto.StringToSignatureAlgorithm(sigAlgo.String())

	servicePrivateKeyHex := privateKey
	_privateKey, err := crypto.DecodePrivateKeyHex(_sigAlgo, servicePrivateKeyHex)
	if err != nil {
		panic(err)
	}
	return _privateKey.PublicKey()
}

func GetReferenceBlockId(flowClient *client.Client) flow.Identifier {
	block, err := flowClient.GetLatestBlock(context.Background(), false)
	if err != nil {
		panic(err)
	}

	return block.ID
}

func CadenceValueToJsonString(value cadence.Value) string {
	result := CadenceValueToInterface(value)
	json1, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}

	return string(json1)
}

func CadenceValueToInterface(field cadence.Value) interface{} {
	switch field.(type) {
	case cadence.Dictionary:
		result := map[string]interface{}{}
		for _, item := range field.(cadence.Dictionary).Pairs {
			result[item.Key.String()] = CadenceValueToInterface(item.Value)
		}
		return result
	case cadence.Struct:
		result := map[string]interface{}{}
		subStructNames := field.(cadence.Struct).StructType.Fields
		for j, subField := range field.(cadence.Struct).Fields {
			result[subStructNames[j].Identifier] = CadenceValueToInterface(subField)
		}
		return result
	case cadence.Array:
		result := []interface{}{}
		for _, item := range field.(cadence.Array).Values {
			result = append(result, CadenceValueToInterface(item))
		}
		return result
	default:
		result, err := strconv.Unquote(field.String())
		if err != nil {
			return field.String()
		}
		return result
	}
}
