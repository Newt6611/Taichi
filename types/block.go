package types

type BlockHeader struct {
	Slot          uint64
	Hash          []byte
	Height        uint64
}

type BlockBody struct {
	Txs []Tx
}

type Block struct {
	Header BlockHeader
	Body   BlockBody
}
