# Taichi

## This is for my personal use only for now.

<div align="center">
    <img src="./assets/logo.png" alt="taichi logo" width="480">
</div>

```
package main

import (
	"fmt"

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

	tai.OnAddressDetach("addr_test1vqushext8jzzns0x9sm6ne2p0m3j3mz0wuu2q0f7hggxjscggvd5t", func (utxo types.UTxO) {
		fmt.Println("dt", utxo.TxHash)
	})

	tai.OnAddressAttach("addr_test1wz4h6068hs93n8j5ar88fgzz6sfnw8krng09xx0mmf36m8c7j9yap", func (utxo types.UTxO) {
		fmt.Println("at", utxo.TxHash)
	})


	tai.Run(93011205, "32af702bc8604616dd2f243a57d6c9a543b2bdba1abdf9eb8b93d0529167585b")
}
```
