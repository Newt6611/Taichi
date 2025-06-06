package taichi

import (
	"github.com/newt6611/taichi/provider"
	"github.com/newt6611/taichi/types"
)

type UTxOHandler func(types.UTxO)

type Taichi struct {
	provider provider.Provider

	addressAttachEvents map[string]UTxOHandler
	addressDetachEvents map[string]UTxOHandler

	paymentAttachEvents map[string]UTxOHandler
	paymentDetachEvents map[string]UTxOHandler

	stakingKeyAttachEvents map[string]UTxOHandler
	stakingKeyDetachEvents map[string]UTxOHandler
}

func NewTaichi(provider provider.Provider) *Taichi {
	return &Taichi{
		provider:             provider,
		addressAttachEvents:  make(map[string]UTxOHandler),
		addressDetachEvents:  make(map[string]UTxOHandler),
		paymentAttachEvents:  make(map[string]UTxOHandler),
		paymentDetachEvents:  make(map[string]UTxOHandler),
		stakingKeyAttachEvents: make(map[string]UTxOHandler),
		stakingKeyDetachEvents: make(map[string]UTxOHandler),
	}
}

func (t *Taichi) Run(slotNum int64, blockHash string) {
	ch := t.provider.Run(slotNum, blockHash)
	for block := range ch {
		for _, tx := range block.Body.Txs {
			for _, input := range tx.Inputs {
				addressStr := input.Address.String()
				if t.addressDetachEvents[addressStr] != nil {
					t.addressDetachEvents[addressStr](input)
				}

				paymentStr := input.Address.PaymentKeyHash().String()
				if t.paymentDetachEvents[paymentStr] != nil {
					t.paymentDetachEvents[paymentStr](input)
				}

				stakeKeyStr := input.Address.StakeKeyHash().String()
				if t.stakingKeyDetachEvents[stakeKeyStr] != nil {
					t.stakingKeyDetachEvents[stakeKeyStr](input)
				}
			}

			for _, output := range tx.Inputs {
				addressStr := output.Address.String()
				if t.addressAttachEvents[addressStr] != nil {
					t.addressAttachEvents[addressStr](output)
				}

				paymentStr := output.Address.PaymentKeyHash().String()
				if t.paymentAttachEvents[paymentStr] != nil {
					t.paymentAttachEvents[paymentStr](output)
				}

				stakeKeyStr := output.Address.StakeKeyHash().String()
				if t.stakingKeyAttachEvents[stakeKeyStr] != nil {
					t.stakingKeyAttachEvents[stakeKeyStr](output)
				}
			}
		}
	}
}

func (t *Taichi) OnAddressAttach(address string, handler UTxOHandler) {
	if _, ok := t.addressAttachEvents[address]; !ok {
		t.addressAttachEvents[address] = handler
	}
}

func (t *Taichi) OnAddressDetach(address string, handler UTxOHandler) {
	if _, ok := t.addressDetachEvents[address]; !ok {
		t.addressDetachEvents[address] = handler
	}
}

func (t *Taichi) OnPaymentKeyHashAttach(paymentKeyHash string, handler UTxOHandler) {
	if _, ok := t.paymentAttachEvents[paymentKeyHash]; !ok {
		t.paymentAttachEvents[paymentKeyHash] = handler
	}
}

func (t *Taichi) OnPaymentKeyHashDetach(paymentKeyHash string, handler UTxOHandler) {
	if _, ok := t.paymentDetachEvents[paymentKeyHash]; !ok {
		t.paymentDetachEvents[paymentKeyHash] = handler
	}
}

func (t *Taichi) OnStakingKeyHashAttach(stakeKeyHash string, handler UTxOHandler) {
	if _, ok := t.stakingKeyAttachEvents[stakeKeyHash]; !ok {
		t.stakingKeyAttachEvents[stakeKeyHash] = handler
	}
}

func (t *Taichi) OnStakingKeyHashDetach(stakeKeyHash string, handler UTxOHandler) {
	if _, ok := t.stakingKeyDetachEvents[stakeKeyHash]; !ok {
		t.stakingKeyDetachEvents[stakeKeyHash] = handler
	}
}
