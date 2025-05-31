package types

type MultiAsset struct {
	PolicyId []byte
	Assets   []Asset
}

type Asset struct {
	Name     []byte
	Quantity int64
}
