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
	PrivateKey              string
	EthClientUrl            string
	Port                    int
	DataDirectory           string
	SupportedEntryPoints    []common.Address
	MaxVerificationGas      *big.Int
	MaxOpsForUnstakedSender int
	Beneficiary             string

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
	viper.SetDefault("erc4337_bundler_supported_entry_points", "0x0F46c65C17AA6b4102046935F33301f0510B163A")
	viper.SetDefault("erc4337_bundler_max_verification_gas", 1500000)
	viper.SetDefault("erc4337_bundler_max_ops_for_unstaked_sender", 4)
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
	_ = viper.BindEnv("erc4337_bundler_eth_client_url")
	_ = viper.BindEnv("erc4337_bundler_private_key")
	_ = viper.BindEnv("erc4337_bundler_port")
	_ = viper.BindEnv("erc4337_bundler_data_directory")
	_ = viper.BindEnv("erc4337_bundler_supported_entry_points")
	_ = viper.BindEnv("erc4337_bundler_beneficiary")
	_ = viper.BindEnv("erc4337_bundler_max_verification_gas")
	_ = viper.BindEnv("erc4337_bundler_max_ops_for_unstaked_sender")
	_ = viper.BindEnv("erc4337_bundler_gin_mode")

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
	maxOpsForUnstakedSender := viper.GetInt("erc4337_bundler_max_ops_for_unstaked_sender")
	debugMode := viper.GetBool("erc4337_bundler_debug_mode")
	ginMode := viper.GetString("erc4337_bundler_gin_mode")
	return &Values{
		PrivateKey:              privateKey,
		EthClientUrl:            ethClientUrl,
		Port:                    port,
		DataDirectory:           dataDirectory,
		SupportedEntryPoints:    supportedEntryPoints,
		Beneficiary:             beneficiary,
		MaxVerificationGas:      maxVerificationGas,
		MaxOpsForUnstakedSender: maxOpsForUnstakedSender,
		DebugMode:               debugMode,
		GinMode:                 ginMode,
		BundlerCollectorTracer:  bct,
	}
}
