package types

import "github.com/blinklabs-io/gouroboros/ledger/common"

type UTxO struct {
	common.Utxo
	TxHash    common.Blake2b256
	Index     uint64
	Address   common.Address
	Assets    []MultiAsset
	DatumHash []byte
	Datum     []byte
	Script    []byte
}
