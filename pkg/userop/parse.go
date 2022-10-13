package userop

import (
	"encoding/hex"
	"errors"
	"math/big"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

type UserOperation struct {
	Sender               string  `json:"sender" mapstructure:"sender" validate:"required,eth_addr"`
	Nonce                big.Int `json:"nonce" mapstructure:"nonce" validate:"required"`
	InitCode             []byte  `json:"initCode"  mapstructure:"initCode" validate:"required"`
	CallData             []byte  `json:"callData" mapstructure:"callData" validate:"required"`
	CallGasLimit         big.Int `json:"callGasLimit" mapstructure:"callGasLimit" validate:"required"`
	VerificationGasLimit big.Int `json:"verificationGasLimit" mapstructure:"verificationGasLimit" validate:"required"`
	PreVerificationGas   big.Int `json:"preVerificationGas" mapstructure:"preVerificationGas" validate:"required"`
	MaxFeePerGas         big.Int `json:"maxFeePerGas" mapstructure:"maxFeePerGas" validate:"required"`
	MaxPriorityFeePerGas big.Int `json:"maxPriorityFeePerGas" mapstructure:"maxPriorityFeePerGas" validate:"required"`
	PaymasterAndData     []byte  `json:"paymasterAndData" mapstructure:"paymasterAndData" validate:"required"`
	Signature            []byte  `json:"signature" mapstructure:"signature" validate:"required"`
}

var validate *validator.Validate

func exactFieldMatch(mapKey, fieldName string) bool {
	return mapKey == fieldName
}

func decodeOpTypes(
	f reflect.Kind,
	t reflect.Kind,
	data interface{}) (interface{}, error) {
	// String to big.Int conversion
	if f == reflect.String && t == reflect.Struct {
		n := new(big.Int)
		n, ok := n.SetString(data.(string), 0)
		if !ok {
			return nil, errors.New("bigInt conversion failed")
		}
		return n, nil
	}

	// Float64 to big.Int conversion
	if f == reflect.Float64 && t == reflect.Struct {
		n, ok := data.(float64)
		if !ok {
			return nil, errors.New("bigInt conversion failed")
		}
		return big.NewInt(int64(n)), nil
	}

	// String to []byte conversion
	if f == reflect.String && t == reflect.Slice {
		byteStr := data.(string)
		if len(byteStr) < 2 || byteStr[:2] != "0x" {
			return nil, errors.New("not byte string")
		}

		b, err := hex.DecodeString(byteStr[2:])
		if err != nil {
			return nil, err
		}
		return b, nil
	}

	return data, nil
}

func validateBigIntType(field reflect.Value) interface{} {
	value, ok := field.Interface().(big.Int)
	if !ok || value.Cmp(big.NewInt(0)) == -1 {
		return nil
	}

	return field
}

func FromMap(data map[string]interface{}) (*UserOperation, error) {
	var op UserOperation

	// Convert map to struct
	config := &mapstructure.DecoderConfig{
		DecodeHook: decodeOpTypes,
		Result:     &op,
		ErrorUnset: true,
		MatchName:  exactFieldMatch,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return nil, err
	}
	if err := decoder.Decode(data); err != nil {
		return nil, err
	}

	// Validate struct
	validate = validator.New()
	validate.RegisterCustomTypeFunc(validateBigIntType, big.Int{})
	err = validate.Struct(op)
	if err != nil {
		return nil, err
	}

	return &op, nil
}
