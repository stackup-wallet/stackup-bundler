package checks

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type codeHash struct {
	Address common.Address `json:"address"`
	Hash    common.Hash    `json:"hash"`
}

func getCodeHashes(ic []common.Address, gc GetCodeFunc) ([]codeHash, error) {
	ret := []codeHash{}

	for _, addr := range ic {
		bytecode, err := gc(addr)
		if err != nil {
			return ret, err
		}

		ret = append(ret, codeHash{
			Address: addr,
			Hash:    crypto.Keccak256Hash(bytecode),
		})
	}

	return ret, nil
}

func hasCodeHashChanges(chs []codeHash, gc GetCodeFunc) (bool, error) {
	prev := map[common.Address]common.Hash{}
	ic := []common.Address{}
	for _, ch := range chs {
		prev[ch.Address] = ch.Hash
		ic = append(ic, ch.Address)
	}

	curr, err := getCodeHashes(ic, gc)
	if err != nil {
		return false, err
	}

	for _, ch := range curr {
		if ch.Hash != prev[ch.Address] {
			return true, nil
		}
	}
	return false, nil
}
