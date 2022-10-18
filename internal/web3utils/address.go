package web3utils

import (
	"reflect"

	"github.com/ethereum/go-ethereum/common"
)

func IsZeroAddress(addr interface{}) bool {
	var address common.Address
	switch v := addr.(type) {
	case string:
		address = common.HexToAddress(v)
	case common.Address:
		address = v
	default:
		return false
	}

	zeroAddressBytes := common.FromHex("0x0000000000000000000000000000000000000000")
	addressBytes := address.Bytes()
	return reflect.DeepEqual(addressBytes, zeroAddressBytes)
}
