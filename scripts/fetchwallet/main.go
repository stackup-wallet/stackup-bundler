package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/viper"
	"github.com/stackup-wallet/stackup-bundler/pkg/signer"
)

func main() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	s, err := signer.New(viper.GetString("erc4337_bundler_private_key"))
	if err != nil {
		panic(fmt.Errorf("fatal signer error: %w", err))
	}
	fmt.Printf("Public key: %s\n", hexutil.Encode(crypto.FromECDSAPub(s.PublicKey))[4:])
	fmt.Printf("Address: %s\n", s.Address)
}
