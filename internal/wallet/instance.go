package wallet

import (
	"crypto/ecdsa"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type Instance struct {
	PrivateKey string
	PublicKey  string
	Address    string
}

func New(pk string) Instance {
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	return Instance{
		PrivateKey: pk,
		PublicKey:  hexutil.Encode(crypto.FromECDSAPub(publicKeyECDSA))[4:],
		Address:    crypto.PubkeyToAddress(*publicKeyECDSA).Hex(),
	}
}
