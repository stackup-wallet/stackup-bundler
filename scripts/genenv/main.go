// Use this for generating a new private key saved to .privatekey
// Implementation from https://goethereumbook.org/en/wallet-generate/
package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/viper"
)

func genPrivateKey() string {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)

	return hexutil.Encode(privateKeyBytes)[2:]
}

func main() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.Set("ERC4337_BUNDLER_ETH_CLIENT_URL", "")
	viper.Set("ERC4337_BUNDLER_PRIVATE_KEY", genPrivateKey())

	if err := viper.WriteConfigAs(".env"); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}
