package BlockFrostChainContext

import (
	"Salvionied/apollo/serialization"
	"Salvionied/apollo/serialization/Address"
	"Salvionied/apollo/serialization/Amount"
	"Salvionied/apollo/serialization/Asset"
	"Salvionied/apollo/serialization/AssetName"
	"Salvionied/apollo/serialization/MultiAsset"
	"Salvionied/apollo/serialization/Policy"
	"Salvionied/apollo/serialization/Redeemer"
	"Salvionied/apollo/serialization/Transaction"
	"Salvionied/apollo/serialization/TransactionInput"
	"Salvionied/apollo/serialization/TransactionOutput"
	"Salvionied/apollo/serialization/UTxO"
	"Salvionied/apollo/serialization/Value"
	"Salvionied/apollo/txBuilding/Backend/Base"
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Salvionied/cbor/v2"
)

type BlockFrostChainContext struct {
	client          *http.Client
	_epoch_info     Base.Epoch
	_epoch          int
	_Network        int
	_genesis_param  Base.GenesisParameters
	_protocol_param Base.ProtocolParameters
	_baseUrl        string
	_projectId      string
	ctx             context.Context
}

func NewBlockfrostChainContext(projectId string, network int, baseUrl string) BlockFrostChainContext {
	ctx := context.Background()
	// latest_epochs, err := api.Epoch(ctx)
	// if err != nil {
	// 	log.Fatal(err, "LATEST EPOCH")
	// }

	bfc := BlockFrostChainContext{client: &http.Client{}, _Network: network, _baseUrl: baseUrl, _projectId: projectId, ctx: ctx}
	bfc.Init()
	return bfc
}
func (bfc *BlockFrostChainContext) Init() {
	latest_epochs := bfc.LatestEpoch()
	bfc._epoch_info = latest_epochs
	//Init Genesis
	params := bfc.GenesisParams()
	bfc._genesis_param = params
	//init epoch
	latest_params := bfc.LatestEpochParams()
	bfc._protocol_param = latest_params
}

func (bfc *BlockFrostChainContext) LatestBlock() Base.Block {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/v0/blocks/latest", bfc._baseUrl), nil)
	req.Header.Set("project_id", bfc._projectId)
	res, err := bfc.client.Do(req)
	if err != nil {
		log.Fatal(err, "REQUEST PROTOCOL")
	}
	body, err := ioutil.ReadAll(res.Body)
	var response Base.Block
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal(err, "UNMARSHAL PROTOCOL")
	}
	return response
}

func (bfc *BlockFrostChainContext) LatestEpoch() Base.Epoch {
	res := Base.Epoch{}
	found := CacheGet[Base.Epoch]("latest_epoch", &res)
	timest := time.Time{}
	foundTime := CacheGet[time.Time]("latest_epoch_time", &timest)
	if !found || !foundTime || time.Since(timest) > 5*time.Minute {
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/v0/epochs/latest", bfc._baseUrl), nil)
		req.Header.Set("project_id", bfc._projectId)
		res, err := bfc.client.Do(req)
		if err != nil {
			log.Fatal(err, "REQUEST PROTOCOL")
		}
		body, err := ioutil.ReadAll(res.Body)
		var response Base.Epoch
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Fatal(err, "UNMARSHAL PROTOCOL")
		}
		CacheSet("latest_epoch", response)
		now := time.Now()
		CacheSet("latest_epoch_time", now)
		return response
	} else {
		return res
	}
}
func (bfc *BlockFrostChainContext) AddressUtxos(address string, gather bool) []Base.AddressUTXO {
	if gather {
		var i = 1
		result := make([]Base.AddressUTXO, 0)
		for {
			req, _ := http.NewRequest("GET", fmt.Sprintf("%s/v0/addresses/%s/utxos?page=%s", bfc._baseUrl, address, fmt.Sprint(i)), nil)
			req.Header.Set("project_id", bfc._projectId)
			res, err := bfc.client.Do(req)
			if err != nil {
				log.Fatal(err, "REQUEST PROTOCOL")
			}
			body, err := ioutil.ReadAll(res.Body)
			var response []Base.AddressUTXO
			err = json.Unmarshal(body, &response)
			if len(response) == 0 {
				break
			}
			if err != nil {
				log.Fatal(err, "UNMARSHAL PROTOCOL")
			}
			result = append(result, response...)
			i++
		}
		return result
	} else {
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/v0/addresses/%s/utxos", bfc._baseUrl, address), nil)
		req.Header.Set("project_id", bfc._projectId)
		res, err := bfc.client.Do(req)
		if err != nil {
			log.Fatal(err, "REQUEST PROTOCOL")
		}
		body, err := ioutil.ReadAll(res.Body)
		var response []Base.AddressUTXO
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Fatal(err, "UNMARSHAL PROTOCOL")
		}
		return response
	}
}

func (bfc *BlockFrostChainContext) LatestEpochParams() Base.ProtocolParameters {
	pm := Base.ProtocolParameters{}
	found := CacheGet[Base.ProtocolParameters]("latest_epoch_params", &pm)
	timest := time.Time{}
	foundTime := CacheGet[time.Time]("latest_epoch_params_time", &timest)
	if !found || !foundTime || time.Since(timest) > 5*time.Minute {
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/v0/epochs/latest/parameters", bfc._baseUrl), nil)
		req.Header.Set("project_id", bfc._projectId)
		res, err := bfc.client.Do(req)
		if err != nil {
			log.Fatal(err, "REQUEST PROTOCOL")
		}
		body, err := ioutil.ReadAll(res.Body)
		var response = Base.ProtocolParameters{}
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Fatal(err, "UNMARSHAL PROTOCOL")
		}
		CacheSet("latest_epoch_params", response)
		now := time.Now()
		CacheSet("latest_epoch_params_time", now)
		return response
	} else {
		return pm
	}
}

func (bfc *BlockFrostChainContext) GenesisParams() Base.GenesisParameters {
	gp := Base.GenesisParameters{}
	found := CacheGet[Base.GenesisParameters]("genesis_params", &gp)
	timestamp := ""
	foundTime := CacheGet[string]("genesis_params_time", &timestamp)
	timest := time.Time{}
	if timestamp != "" {
		timest, _ = time.Parse(time.RFC3339, timestamp)
	}
	if !found || !foundTime || time.Since(timest) > 5*time.Minute {
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/v0/Genesis", bfc._baseUrl), nil)
		req.Header.Set("project_id", bfc._projectId)
		res, err := bfc.client.Do(req)
		if err != nil {
			log.Fatal(err, "REQUEST PROTOCOL")
		}
		body, err := ioutil.ReadAll(res.Body)
		var response = Base.GenesisParameters{}
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Fatal(err, "UNMARSHAL PROTOCOL")
		}
		CacheSet("genesis_params", response)
		now := time.Now()
		CacheSet("genesis_params_time", now)
		return response
	} else {

		return gp
	}
}
func (bfc *BlockFrostChainContext) _CheckEpochAndUpdate() bool {
	if bfc._epoch_info.EndTime <= int(time.Now().Unix()) {
		latest_epochs := bfc.LatestEpoch()
		bfc._epoch_info = latest_epochs
		return true
	}
	return false
}

func (bfc *BlockFrostChainContext) Network() int {
	return bfc._Network
}

func (bfc *BlockFrostChainContext) Epoch() int {
	if bfc._CheckEpochAndUpdate() {
		new_epoch := bfc.LatestEpoch()
		bfc._epoch = new_epoch.Epoch
	}
	return bfc._epoch
}

func (bfc *BlockFrostChainContext) LastBlockSlot() int {
	block := bfc.LatestBlock()
	return block.Slot
}

func (bfc *BlockFrostChainContext) GetGenesisParams() Base.GenesisParameters {
	if bfc._CheckEpochAndUpdate() {
		params := bfc.GenesisParams()
		bfc._genesis_param = params
	}
	return bfc._genesis_param
}

func (bfc *BlockFrostChainContext) GetProtocolParams() Base.ProtocolParameters {
	if bfc._CheckEpochAndUpdate() {
		latest_params := bfc.LatestEpochParams()
		bfc._protocol_param = latest_params
	}
	return bfc._protocol_param
}

func (bfc *BlockFrostChainContext) MaxTxFee() int {
	protocol_param := bfc.GetProtocolParams()
	maxTxExSteps, _ := strconv.Atoi(protocol_param.MaxTxExSteps)
	maxTxExMem, _ := strconv.Atoi(protocol_param.MaxTxExMem)
	return Base.Fee(bfc, protocol_param.MaxTxSize, maxTxExSteps, maxTxExMem)
}

func (bfc *BlockFrostChainContext) Utxos(address Address.Address) []UTxO.UTxO {
	results := bfc.AddressUtxos(address.String(), true)
	utxos := make([]UTxO.UTxO, 0)
	for _, result := range results {
		decodedTxId, _ := hex.DecodeString(result.TxHash)
		tx_in := TransactionInput.TransactionInput{TransactionId: decodedTxId, Index: result.OutputIndex}
		amount := result.Amount
		lovelace_amount := 0
		multi_assets := MultiAsset.MultiAsset[int64]{}
		for _, item := range amount {
			if item.Unit == "lovelace" {
				amount, err := strconv.Atoi(item.Quantity)
				if err != nil {
					log.Fatal(err)
				}
				lovelace_amount += amount
			} else {
				asset_quantity, err := strconv.ParseInt(item.Quantity, 10, 64)
				if err != nil {
					log.Fatal(err)
				}
				policy_id := Policy.PolicyId{Value: item.Unit[:56]}
				asset_name := *AssetName.NewAssetNameFromHexString(item.Unit[56:])
				_, ok := multi_assets[policy_id]
				if !ok {
					multi_assets[policy_id] = Asset.Asset[int64]{}
				}
				multi_assets[policy_id][asset_name] = int64(asset_quantity)
			}
		}
		final_amount := Value.Value{}
		if len(multi_assets) > 0 {
			final_amount = Value.Value{Am: Amount.Amount{Coin: int64(lovelace_amount), Value: multi_assets}, HasAssets: true}
		} else {
			final_amount = Value.Value{Coin: int64(lovelace_amount), HasAssets: false}
		}
		datum_hash := serialization.DatumHash{}
		if result.DataHash != "" && result.InlineDatum == "" {

			datum_hash = serialization.DatumHash{}
			copy(datum_hash.Payload[:], result.DataHash[:])
		}
		// if result.InlineDatum != "" {
		// 	decoded, err := hex.DecodeString(result.InlineDatum)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// 	var x serialization.PlutusData
		// 	err = cbor.Unmarshal(decoded, &x)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// 	datum := x
		// }
		tx_out := TransactionOutput.TransactionOutput{PreAlonzo: TransactionOutput.TransactionOutputShelley{
			Address:   address,
			Amount:    final_amount,
			DatumHash: datum_hash,
			HasDatum:  len(datum_hash.Payload) > 0}, IsPostAlonzo: false}
		utxos = append(utxos, UTxO.UTxO{Input: tx_in, Output: tx_out})
	}
	return utxos
}

func (bfc *BlockFrostChainContext) SubmitTx(tx Transaction.Transaction) serialization.TransactionId {
	txBytes, _ := cbor.Marshal(tx)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/v0/tx/submit", bfc._baseUrl), bytes.NewBuffer(txBytes))
	req.Header.Set("project_id", bfc._projectId)
	req.Header.Set("Content-Type", "application/cbor")
	res, err := bfc.client.Do(req)
	if err != nil {
		log.Fatal(err, "REQUEST PROTOCOL")
	}
	body, err := ioutil.ReadAll(res.Body)
	var response any
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal(err, "UNMARSHAL PROTOCOL")
	}
	return serialization.TransactionId{Payload: tx.TransactionBody.Hash()}
}

type EvalResult struct {
	Result map[string]map[string]int `json:"EvaluationResult"`
}

type ExecutionResult struct {
	Result EvalResult `json:"result"`
}

func (bfc *BlockFrostChainContext) EvaluateTx(tx []byte) map[string]Redeemer.ExecutionUnits {
	encoded := hex.EncodeToString(tx)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/v0/utils/txs/evaluate", bfc._baseUrl), strings.NewReader(encoded))
	req.Header.Set("project_id", bfc._projectId)
	req.Header.Set("Content-Type", "application/cbor")
	res, err := bfc.client.Do(req)
	if err != nil {
		log.Fatal(err, "REQUEST PROTOCOL")
	}
	body, err := ioutil.ReadAll(res.Body)
	var x any
	err = json.Unmarshal(body, &x)
	if err != nil {
		log.Fatal(err, "UNMARSHAL PROTOCOL")
	}
	var response ExecutionResult
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal(err, "UNMARSHAL PROTOCOL")
	}
	final_result := make(map[string]Redeemer.ExecutionUnits, 0)
	for k, v := range response.Result.Result {

		final_result[k] = Redeemer.ExecutionUnits{Steps: int64(v["steps"]), Mem: int64(v["memory"])}
	}
	return final_result
}