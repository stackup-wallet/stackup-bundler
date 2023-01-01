package config

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stackup-wallet/stackup-bundler/pkg/signer"
	"github.com/stackup-wallet/stackup-bundler/pkg/tracer"
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
	DebugMode              bool
	GinMode                string
	BundlerCollectorTracer string
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
	viper.SetDefault("erc4337_bundler_supported_entry_points", "0x1306b01bC3e4AD202612D3843387e94737673F53")
	viper.SetDefault("erc4337_bundler_max_verification_gas", 1500000)
	viper.SetDefault("erc4337_bundler_debug_mode", false)
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
	if err := viper.BindEnv("erc4337_bundler_eth_client_url"); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("erc4337_bundler_private_key"); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("erc4337_bundler_port"); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("erc4337_bundler_data_directory"); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("erc4337_bundler_supported_entry_points"); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("erc4337_bundler_beneficiary"); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("erc4337_bundler_max_verification_gas"); err != nil {
		panic(err)
	}
	if err := viper.BindEnv("erc4337_bundler_gin_mode"); err != nil {
		panic(err)
	}

	// Validate required variables
	if !viper.IsSet("erc4337_bundler_eth_client_url") ||
		viper.GetString("erc4337_bundler_eth_client_url") == "" {
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

	// Load js tracer from embedded file
	bct, err := tracer.Load()
	if err != nil {
		panic(err)
	}

	// Return Values
	privateKey := viper.GetString("erc4337_bundler_private_key")
	ethClientUrl := viper.GetString("erc4337_bundler_eth_client_url")
	port := viper.GetInt("erc4337_bundler_port")
	dataDirectory := viper.GetString("erc4337_bundler_data_directory")
	supportedEntryPoints := envArrayToAddressSlice(viper.GetString("erc4337_bundler_supported_entry_points"))
	beneficiary := viper.GetString("erc4337_bundler_beneficiary")
	maxVerificationGas := big.NewInt(int64(viper.GetInt("erc4337_bundler_max_verification_gas")))
	debugMode := viper.GetBool("erc4337_bundler_debug_mode")
	ginMode := viper.GetString("erc4337_bundler_gin_mode")
	return &Values{
		PrivateKey:             privateKey,
		EthClientUrl:           ethClientUrl,
		Port:                   port,
		DataDirectory:          dataDirectory,
		SupportedEntryPoints:   supportedEntryPoints,
		Beneficiary:            beneficiary,
		MaxVerificationGas:     maxVerificationGas,
		DebugMode:              debugMode,
		GinMode:                ginMode,
		BundlerCollectorTracer: bct,
	}
}
