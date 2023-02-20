package checks

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/sync/errgroup"
)

type codeHash struct {
	Address common.Address `json:"address"`
	Hash    common.Hash    `json:"hash"`
}

func getCodeHashAsync(addr common.Address, gc GetCodeFunc, c chan codeHash) func() error {
	return func() error {
		bytecode, err := gc(addr)
		if err != nil {
			c <- codeHash{}
			return err
		}

		ch := codeHash{
			Address: addr,
			Hash:    crypto.Keccak256Hash(bytecode),
		}
		c <- ch
		return nil
	}
}

func getCodeHashes(ic []common.Address, gc GetCodeFunc) ([]codeHash, error) {
	g := new(errgroup.Group)
	c := make(chan codeHash)
	ret := []codeHash{}

	for _, addr := range ic {
		g.Go(getCodeHashAsync(addr, gc, c))
	}
	for range ic {
		ret = append(ret, <-c)
	}
	if err := g.Wait(); err != nil {
		return ret, err
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
