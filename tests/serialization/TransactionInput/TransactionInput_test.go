package TransactionInput

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/salvionied/apollo/serialization/TransactionInput"
)

func TestJsonMarshal(t *testing.T) {
	var testDefs = []struct {
		txId         string
		index        int
		expectedJson string
	}{
		{
			txId:         "d2153af861591c5cfe039de304f1e408edbf8bbfc7854621625bb74a4f6cd5cb",
			index:        0,
			expectedJson: `"d2153af861591c5cfe039de304f1e408edbf8bbfc7854621625bb74a4f6cd5cb.0"`,
		},
		{
			txId:         "e1ce40b9684c1d074be2e4b0c8abd5ccbc33ab6384bd214b032279247f7bb470",
			index:        1,
			expectedJson: `"e1ce40b9684c1d074be2e4b0c8abd5ccbc33ab6384bd214b032279247f7bb470.1"`,
		},
	}
	for _, test := range testDefs {
		txIdBytes, err := hex.DecodeString(test.txId)
		if err != nil {
			t.Fatalf("unexpected failure decoding hex: %s", err)
		}
		tmpObj := TransactionInput.TransactionInput{
			TransactionId: txIdBytes,
			Index:         test.index,
		}
		tmpJson, err := json.Marshal(&tmpObj)
		if err != nil {
			t.Fatalf("unexpected failure marshaling to JSON: %s", err)
		}
		if string(tmpJson) != test.expectedJson {
			t.Fatalf("object did not marshal to expected JSON\n  got: %s\n  wanted: %s", tmpJson, test.expectedJson)
		}
	}
}
