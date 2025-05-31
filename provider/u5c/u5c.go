package u5c

import (
	"fmt"

	"github.com/blinklabs-io/gouroboros/ledger/common"
	"github.com/newt6611/taichi/types"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/cardano"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/sync"
	utxorpc "github.com/utxorpc/go-sdk"
)

type U5C struct {
	client *utxorpc.UtxorpcClient
}

func NewU5C(baseUrl string, header map[string]string) *U5C {
	client := utxorpc.NewClient(utxorpc.WithBaseUrl(baseUrl))
	for k, v := range header {
		client.SetHeader(k, v)
	}
	return &U5C{
		client: client,
	}
}

func (u5c *U5C) Run(slotNum int64, blockHash string) <-chan types.Block {
	resultCh := make(chan types.Block, 1)

	stream, err := u5c.client.FollowTip(blockHash, slotNum)
	if err != nil {
		utxorpc.HandleError(err)
		return nil
	}

	go func() {
		for stream.Receive() {
			resp := stream.Msg()
			action := resp.GetAction()
			switch action := action.(type) {
			case *sync.FollowTipResponse_Apply:
				block := parseBlock(action.Apply.GetCardano())
				resultCh <- block

			case *sync.FollowTipResponse_Undo:
				fmt.Println("Action: Undo")
			case *sync.FollowTipResponse_Reset_:
				fmt.Println("Action: Reset")
			default:
				fmt.Println("Unknown action type")
			}
		}
	}()

	return resultCh
}

func parseBlock(block *cardano.Block) types.Block {
	resultBlock := types.Block{}
	// Header
	resultBlock.Header = types.BlockHeader{
		Slot:   block.Header.Slot,
		Hash:   block.Header.Hash,
		Height: block.Header.Height,
	}

	// Body
	txs := block.Body.GetTx()
	resultBlock.Body.Txs = make([]types.Tx, len(txs))

	for txIdx, tx := range txs {
		resultBlock.Body.Txs[txIdx].Inputs = make([]types.UTxO, len(tx.Inputs))
		resultBlock.Body.Txs[txIdx].Outputs = make([]types.UTxO, len(tx.Outputs))

		// Inputs
		for inputIdx, txIn := range tx.Inputs {
			out := txIn.AsOutput

			// Decode address
			address := common.Address{}
			err := address.UnmarshalCBOR(out.Address)
			if err != nil {
				fmt.Println("Failed to parse address:", err)
				continue
			}

			// Prepare input UTxO
			utxo := types.UTxO{
				TxHash:    common.NewBlake2b256(txIn.TxHash),
				Index:     uint64(txIn.OutputIndex),
				DatumHash: out.Datum.Hash,
				Datum:     out.Datum.OriginalCbor,
				Address:   address,
				Assets:    make([]types.MultiAsset, len(out.Assets)),
			}

			// Extract assets
			utxo.Assets = parseMultiAssets(out.Assets)

			// Assign back to result block
			resultBlock.Body.Txs[txIdx].Inputs[inputIdx] = utxo
		}
		// Outputs
		for outputIdx, txOut := range tx.Outputs {
			// Decode address
			address := common.Address{}
			err := address.UnmarshalCBOR(txOut.Address)
			// TODO: Better error handling
			if err != nil {
				fmt.Println("Failed to parse address:", err)
			}

			// Prepare input UTxO
			utxo := types.UTxO{
				TxHash:    common.NewBlake2b256(tx.Hash),
				Index:     uint64(outputIdx),
				DatumHash: txOut.Datum.Hash,
				Datum:     txOut.Datum.OriginalCbor,
				Address:   address,
				Assets:    make([]types.MultiAsset, len(txOut.Assets)),
			}

			// Extract assets
			utxo.Assets = parseMultiAssets(txOut.Assets)
			utxo.Address.PaymentKeyHash()

			// Assign back to result block
			resultBlock.Body.Txs[txIdx].Outputs[outputIdx] = utxo
		}
	}

	return resultBlock
}

func parseMultiAssets(rawAssets []*cardano.Multiasset) []types.MultiAsset {
	multiAssets := make([]types.MultiAsset, len(rawAssets))

	for i, ma := range rawAssets {
		assets := make([]types.Asset, len(ma.Assets))
		for j, asset := range ma.Assets {
			assets[j] = types.Asset{
				Name:     asset.Name,
				Quantity: int64(asset.OutputCoin),
			}
		}
		multiAssets[i] = types.MultiAsset{
			PolicyId: ma.PolicyId,
			Assets:   assets,
		}
	}

	return multiAssets
}
