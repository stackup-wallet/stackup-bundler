package config

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stackup-wallet/stackup-bundler/pkg/signer"
)

type Values struct {
	// Documented variables.
	PrivateKey           string
	EthClientUrl         string
	Port                 int
	DataDirectory        string
	SupportedEntryPoints []common.Address
	MaxVerificationGas   *big.Int
	Beneficiary          string

	// Undocumented variables.
	GinMode string
}

func envArrayToAddressSlice(s string) []common.Address {
	env := strings.Split(s, ",")
	slc := []common.Address{}
	for _, ep := range env {
		slc = append(slc, common.HexToAddress(strings.TrimSpace(ep)))
	}

	return slc
}

// GetValues returns config for the bundler that has been read in from env vars.
// See https://docs.stackup.sh/docs/packages/bundler/configure for details.
func GetValues() *Values {
	// Default variables
	viper.SetDefault("erc4337_bundler_port", 4337)
	viper.SetDefault("erc4337_bundler_data_directory", "/tmp/stackup_bundler")
	viper.SetDefault("erc4337_bundler_supported_entry_points", "0x1D9a2CB3638C2FC8bF9C01D088B79E75CD188b17")
	viper.SetDefault("erc4337_bundler_max_verification_gas", 1500000)
	viper.SetDefault("erc4337_bundler_gin_mode", gin.ReleaseMode)

	// Read in from .env file if available
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found
			// Can ignore
		} else {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}

	// Read in from environment variables
	viper.BindEnv("erc4337_bundler_eth_client_url")
	viper.BindEnv("erc4337_bundler_private_key")
	viper.BindEnv("erc4337_bundler_port")
	viper.BindEnv("erc4337_bundler_data_directory")
	viper.BindEnv("erc4337_bundler_supported_entry_points")
	viper.BindEnv("erc4337_bundler_beneficiary")
	viper.BindEnv("erc4337_bundler_max_verification_gas")
	viper.BindEnv("erc4337_bundler_gin_mode")

	// Validate required variables
	if !viper.IsSet("erc4337_bundler_eth_client_url") || viper.GetString("erc4337_bundler_eth_client_url") == "" {
		panic("Fatal config error: erc4337_bundler_eth_client_url not set")
	}

	if !viper.IsSet("erc4337_bundler_private_key") || viper.GetString("erc4337_bundler_private_key") == "" {
		panic("Fatal config error: erc4337_bundler_private_key not set")
	}

	if !viper.IsSet("erc4337_bundler_beneficiary") {
		s, err := signer.New(viper.GetString("erc4337_bundler_private_key"))
		if err != nil {
			panic(err)
		}
		viper.SetDefault("erc4337_bundler_beneficiary", s.Address.String())
	}

	// Return Values
	return &Values{
		PrivateKey:           viper.GetString("erc4337_bundler_private_key"),
		EthClientUrl:         viper.GetString("erc4337_bundler_eth_client_url"),
		Port:                 viper.GetInt("erc4337_bundler_port"),
		DataDirectory:        viper.GetString("erc4337_bundler_data_directory"),
		SupportedEntryPoints: envArrayToAddressSlice(viper.GetString("erc4337_bundler_supported_entry_points")),
		Beneficiary:          viper.GetString("erc4337_bundler_beneficiary"),
		MaxVerificationGas:   big.NewInt(int64(viper.GetInt("erc4337_bundler_max_verification_gas"))),
		GinMode:              viper.GetString("erc4337_bundler_gin_mode"),
	}
}
