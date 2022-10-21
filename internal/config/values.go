package config

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stackup-wallet/stackup-bundler/internal/wallet"
)

type Values struct {
	// Documented variables.
	PrivateKey           string
	RedisUrl             string
	RpcUrl               string
	Port                 int
	SupportedEntryPoints []string
	Beneficiary          string

	// Undocumented variables.
	GinMode string
}

func envArrayToSlice(s string) []string {
	slc := strings.Split(s, ",")
	for i := range slc {
		slc[i] = strings.TrimSpace(slc[i])
	}

	return slc
}

// Returns config values for the RPC server that has been read in from env vars.
// See https://docs.stackup.sh/docs/packages/client/configuration for details.
func GetValues() Values {
	// Default variables
	viper.SetDefault("erc4337_bundler_port", 4337)
	viper.SetDefault("erc4337_bundler_supported_entry_points", "0x1b98F08dB8F12392EAE339674e568fe29929bC47")
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
	viper.BindEnv("erc4337_bundler_private_key")
	viper.BindEnv("erc4337_bundler_redis_url")
	viper.BindEnv("erc4337_bundler_rpc_url")
	viper.BindEnv("erc4337_bundler_port")
	viper.BindEnv("erc4337_bundler_supported_entry_points")
	viper.BindEnv("erc4337_bundler_beneficiary")
	viper.BindEnv("erc4337_bundler_gin_mode")

	// Validate required variables
	if !viper.IsSet("erc4337_bundler_private_key") || viper.GetString("erc4337_bundler_private_key") == "" {
		panic("Fatal config error: erc4337_bundler_private_key not set")
	}

	if !viper.IsSet("erc4337_bundler_redis_url") || viper.GetString("erc4337_bundler_redis_url") == "" {
		panic("Fatal config error: erc4337_bundler_redis_url not set")
	}

	if !viper.IsSet("erc4337_bundler_rpc_url") || viper.GetString("erc4337_bundler_rpc_url") == "" {
		panic("Fatal config error: erc4337_bundler_rpc_url not set")
	}

	if !viper.IsSet("erc4337_bundler_beneficiary") {
		w := wallet.New(viper.GetString("erc4337_bundler_private_key"))
		viper.SetDefault("erc4337_bundler_beneficiary", w.Address)
	}

	// Return Values
	return Values{
		PrivateKey:           viper.GetString("erc4337_bundler_private_key"),
		RedisUrl:             viper.GetString("erc4337_bundler_redis_url"),
		RpcUrl:               viper.GetString("erc4337_bundler_rpc_url"),
		Port:                 viper.GetInt("erc4337_bundler_port"),
		SupportedEntryPoints: envArrayToSlice(viper.GetString("erc4337_bundler_supported_entry_points")),
		Beneficiary:          viper.GetString("erc4337_bundler_beneficiary"),
		GinMode:              viper.GetString("erc4337_bundler_gin_mode"),
	}
}
