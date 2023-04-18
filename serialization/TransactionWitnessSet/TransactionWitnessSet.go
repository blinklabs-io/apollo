package TransactionWitnessSet

import (
	"Salvionied/apollo/serialization/NativeScript"
	"Salvionied/apollo/serialization/PlutusData"
	"Salvionied/apollo/serialization/Redeemer"
	"Salvionied/apollo/serialization/VerificationKeyWitness"
)

type TransactionWitnessSet struct {
	VkeyWitnesses      []VerificationKeyWitness.VerificationKeyWitness `cbor:"0,keyasint,omitempty"`
	NativeScripts      []NativeScript.NativeScript                     `cbor:"1,keyasint,omitempty"`
	BootstrapWitnesses []any                                           `cbor:"2,keyasint,omitempty"`
	PlutusV1Script     []PlutusData.PlutusV1Script                     `cbor:"3,keyasint,omitempty"`
	PlutusV2Script     []PlutusData.PlutusV2Script                     `cbor:"6,keyasint,omitempty"`
	PlutusData         []PlutusData.PlutusData                         `cbor:"4,keyasint,omitempty"`
	Redeemer           []Redeemer.Redeemer                             `cbor:"5,keyasint,omitempty"`
}