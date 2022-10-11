package main

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/stackup-wallet/stackup-bundler/internal/wallet"
)

func main() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	w := wallet.New(viper.GetString("erc4337_bundler_private_key"))
	fmt.Printf("Public key: %s\n", w.PublicKey)
	fmt.Printf("Address: %s\n", w.Address)
}
