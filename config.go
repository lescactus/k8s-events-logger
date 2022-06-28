package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const (
	// Name of the application
	AppName = "k8s-events-logger"
)

type AppConfig struct {
	*viper.Viper
}

func NewAppConfig() *AppConfig {
	config := &AppConfig{
		Viper: viper.New(),
	}

	// Set default configurations
	config.setDefaults()

	// Select the .env file
	config.SetConfigName(config.GetString("APP_CONFIG_NAME"))
	config.SetConfigType("dotenv")
	config.AddConfigPath(config.GetString("APP_CONFIG_PATH"))

	// Automatically refresh environment variables
	config.AutomaticEnv()

	// Read configuration
	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Println("failed to read configuration:", err.Error())
			os.Exit(1)
		}
	}

	if !config.isValidOutput() {
		fmt.Printf("invalid value for OUTPUT: %s. Valid values are 'console' or 'json'", config.GetString("OUTPUT"))
		os.Exit(1)
	}

	return config
}

func (config *AppConfig) setDefaults() {
	// Set default App configuration
	config.SetDefault("APP_CONFIG_NAME", ".env")
	config.SetDefault("APP_CONFIG_PATH", ".")

	config.SetDefault("APP_ADDR", ":8080")

	config.SetDefault("NAMESPACES", "default")
	config.SetDefault("OUTPUT", "console") // Available: "console", "json"
}

func (config *AppConfig) isValidOutput() bool {
	output := config.GetString("OUTPUT")

	if (output != "console") && (output != "json") {
		return false
	}
	return true
}
