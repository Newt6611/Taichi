package provider

import (
	"github.com/newt6611/taichi/types"
)

type Provider interface {
	Run(slotNum int64, blockHash string) <-chan types.Block
}
