package example

import (
	"github.com/newt6611/taichi"
	"github.com/newt6611/taichi/provider/u5c"
	"github.com/newt6611/taichi/types"
)

func main() {
	const url = "utxorpc_url"
	const apiKey = "utxorpc_api_key"
	u5c := u5c.NewU5C(url, map[string]string{
		"dmtr-api-key": apiKey,
	})

	tai := taichi.NewTaichi(u5c)

	address := "addr_test1vqushext8jzzns0x9sm6ne2p0m3j3mz0wuu2q0f7hggxjscggvd5t"

	tai.OnAddressDetach(address, func(utxo types.UTxO) {
		// When an address utxo is spent
	})

	tai.OnAddressAttach(address, func(utxo types.UTxO) {
		// When an address receives a new utxo
	})

	paymentKey := "036cea874e068fd88c786c92628241f267decc992116e0e2cea7299a"

	tai.OnPaymentKeyHashDetach(paymentKey, func(utxo types.UTxO) {
		// When an payment credential utxo is spent
	})

	tai.OnPaymentKeyHashAttach(paymentKey, func(utxo types.UTxO) {
		// When an payment credential receives a new utxo
	})

	stakingKey := "e34fec0560cbb68fbf0e8758d52ec25c56574bfe0a7e332974714f1d"

	tai.OnStakeKeyHashDetach(stakingKey, func(utxo types.UTxO) {
		// When an staking credential utxo is spent
	})

	tai.OnStakeKeyHashAttach(stakingKey, func(utxo types.UTxO) {
		// When an staking credential receives a new utxo
	})

	tai.Run(93011205, "32af702bc8604616dd2f243a57d6c9a543b2bdba1abdf9eb8b93d0529167585b")
}
