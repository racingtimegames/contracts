package transactions

import (
	"context"
	"contracts/utils/lib"
	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
)

type BatchMintInfo struct {
	Name                string   `json:"name,omitempty" gorm:"index;column:name"`                           // 名称
	Artist              string   `json:"artist,omitempty" gorm:"column:artist"`                             // 艺术家
	ArtistIntroduction  string   `json:"artist_introduction,omitempty" gorm:"column:artist_introduction"`   // 艺术家介绍
	ArtworkIntroduction string   `json:"artwork_introduction,omitempty" gorm:"column:artwork_introduction"` // 作品介绍
	TypeId              uint64   `json:"type_id,omitempty" gorm:"index;column:type_id"`                     // NFT类型的ID
	NFTType             string   `json:"nft_type,omitempty" gorm:"column:nft_type"`                         // NFT类型
	Description         string   `json:"description,omitempty" gorm:"column:description"`                   // 描述
	IpfsLink            string   `json:"ipfs_link,omitempty" gorm:"column:ipfs_link"`                       // ipfs链接
	MD5Hash             string   `json:"md5_hash,omitempty" gorm:"column:md5_hash"`                         // 图片的md5哈希
	TotalNumber         uint32   `json:"total_number,omitempty" gorm:"column:total_number"`                 // NFT创建数量
	Address             string   `json:"address" gorm:"column:address"`                                     // 接受地址
	SerialNumber        []uint32 `json:"serial_number" gorm:"foreignKey:ID;references:ID"`                  // mint序号
}

func BatchMintNFT(client *client.Client, script string, info BatchMintInfo, auth lib.Account) (tx *flow.Transaction, err error) {

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
		SetGasLimit(9999).
		SetProposalKey(acctAddress, acctKey.Index, acctKey.SequenceNumber).
		SetReferenceBlockID(referenceBlock.ID).
		SetPayer(acctAddress).
		AddAuthorizer(acctAddress)

	if err = tx.AddArgument(cadence.NewAddress(flow.HexToAddress(info.Address))); err != nil {
		return
	}
	if err = tx.AddArgument(cadence.String(info.Name)); err != nil {
		return
	}
	if err = tx.AddArgument(cadence.String(info.Artist)); err != nil {
		return
	}
	if err = tx.AddArgument(cadence.String(info.ArtistIntroduction)); err != nil {
		return
	}
	if err = tx.AddArgument(cadence.String(info.ArtworkIntroduction)); err != nil {
		return
	}
	if err = tx.AddArgument(cadence.NewUInt64(info.TypeId)); err != nil {
		return
	}
	if err = tx.AddArgument(cadence.String(info.NFTType)); err != nil {
		return
	}
	if err = tx.AddArgument(cadence.String(info.Description)); err != nil {
		return
	}
	if err = tx.AddArgument(cadence.String(info.IpfsLink)); err != nil {
		return
	}
	if err = tx.AddArgument(cadence.String(info.MD5Hash)); err != nil {
		return
	}
	var values []cadence.Value
	for _, id := range info.SerialNumber {
		values = append(values, cadence.NewUInt32(id))
	}

	if err = tx.AddArgument(cadence.NewArray(values)); err != nil {
		return
	}
	if err = tx.AddArgument(cadence.NewUInt32(info.TotalNumber)); err != nil {
		return
	}

	if err = tx.SignEnvelope(acctAddress, acctKey.Index, signer); err != nil {
		return
	}

	if err = client.SendTransaction(context.Background(), *tx); err != nil {
		return
	}
	//result, err = lib.WaitForSeal(context.Background(), global.FlowClient, tx.ID())
	return
}
