package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	BMCType BMCType     `yaml:"bmc_type" mapstructure:"bmc_type"`
	ILO     ILOConfig   `yaml:"ilo" mapstructure:"ilo"`
	IDRAC   IDRACConfig `yaml:"idrac" mapstructure:"idrac"`
}

// ILOConfig represents iLO connection configuration
type ILOConfig struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Username string `yaml:"username" mapstructure:"username"`
	Password string `yaml:"password" mapstructure:"password"`
	Port     int    `yaml:"port" mapstructure:"port"`
	UseHTTPS bool   `yaml:"use_https" mapstructure:"use_https"`
}

// IDRACConfig represents iDRAC connection configuration
type IDRACConfig struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Username string `yaml:"username" mapstructure:"username"`
	Password string `yaml:"password" mapstructure:"password"`
	Port     int    `yaml:"port" mapstructure:"port"`
	UseHTTPS bool   `yaml:"use_https" mapstructure:"use_https"`
}

var config Config

func loadConfig() error {
	// Set default values
	viper.SetDefault("bmc_type", "ilo")
	viper.SetDefault("ilo.port", 443)
	viper.SetDefault("ilo.use_https", true)
	viper.SetDefault("idrac.port", 443)
	viper.SetDefault("idrac.use_https", true)

	// Environment variable bindings
	viper.AutomaticEnv()

	// Bind BMC type environment variable
	viper.BindEnv("bmc_type", "BMC_TYPE")

	// Bind iLO specific environment variables
	viper.BindEnv("ilo.host", "ILO_HOST")
	viper.BindEnv("ilo.username", "ILO_USERNAME")
	viper.BindEnv("ilo.password", "ILO_PASSWORD")
	viper.BindEnv("ilo.port", "ILO_PORT")
	viper.BindEnv("ilo.use_https", "ILO_USE_HTTPS")

	// Bind iDRAC specific environment variables
	viper.BindEnv("idrac.host", "IDRAC_HOST")
	viper.BindEnv("idrac.username", "IDRAC_USERNAME")
	viper.BindEnv("idrac.password", "IDRAC_PASSWORD")
	viper.BindEnv("idrac.port", "IDRAC_PORT")
	viper.BindEnv("idrac.use_https", "IDRAC_USE_HTTPS")

	// Configuration file handling
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Look for config in current directory
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	// Read config file if it exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found is okay, we can work with env vars
		if verbose {
			fmt.Printf("Config file not found, using environment variables and defaults\n")
		}
	} else {
		if verbose {
			fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())
		}
	}

	// Unmarshal the configuration
	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("error unmarshaling config: %w", err)
	}

	return validateConfig()
}

func validateConfig() error {
	switch config.BMCType {
	case BMCTypeILO:
		if config.ILO.Host == "" {
			return fmt.Errorf("iLO host is required (set ILO_HOST environment variable or host in config file)")
		}
		if config.ILO.Username == "" {
			return fmt.Errorf("iLO username is required (set ILO_USERNAME environment variable or username in config file)")
		}
		if config.ILO.Password == "" {
			return fmt.Errorf("iLO password is required (set ILO_PASSWORD environment variable or password in config file)")
		}
	case BMCTypeIDRAC:
		if config.IDRAC.Host == "" {
			return fmt.Errorf("iDRAC host is required (set IDRAC_HOST environment variable or host in config file)")
		}
		if config.IDRAC.Username == "" {
			return fmt.Errorf("iDRAC username is required (set IDRAC_USERNAME environment variable or username in config file)")
		}
		if config.IDRAC.Password == "" {
			return fmt.Errorf("iDRAC password is required (set IDRAC_PASSWORD environment variable or password in config file)")
		}
	default:
		return fmt.Errorf("unsupported BMC type: %s (supported types: ilo, idrac)", config.BMCType)
	}
	return nil
}

func createSampleConfig() error {
	sampleConfig := `# BMC CLI Configuration File
# Specify the BMC type: 'ilo' for HPE iLO or 'idrac' for DELL iDRAC
bmc_type: "ilo"

# HPE iLO Configuration
ilo:
  host: "192.168.1.100"          # iLO IP address or hostname
  username: "admin"              # iLO username
  password: "password"           # iLO password  
  port: 443                      # iLO port (default: 443)
  use_https: true                # Use HTTPS (default: true)

# DELL iDRAC Configuration
idrac:
  host: "192.168.1.101"          # iDRAC IP address or hostname
  username: "root"               # iDRAC username
  password: "calvin"             # iDRAC password
  port: 443                      # iDRAC port (default: 443)
  use_https: true                # Use HTTPS (default: true)
`

	configPath := filepath.Join(".", "config.yaml")
	if _, err := os.Stat(configPath); err == nil {
		return fmt.Errorf("config file already exists at %s", configPath)
	}

	if err := os.WriteFile(configPath, []byte(sampleConfig), 0644); err != nil {
		return fmt.Errorf("error creating sample config: %w", err)
	}

	fmt.Printf("Sample configuration file created at %s\n", configPath)
	fmt.Println("Please edit the file with your BMC credentials and settings.")
	return nil
}

// NewBMCClient creates a BMC client based on the configuration
func NewBMCClient() (BMCClient, error) {
	switch config.BMCType {
	case BMCTypeILO:
		return NewILOClient(
			config.ILO.Host,
			config.ILO.Username,
			config.ILO.Password,
			config.ILO.Port,
			config.ILO.UseHTTPS,
		), nil
	case BMCTypeIDRAC:
		return NewIDRACClient(
			config.IDRAC.Host,
			config.IDRAC.Username,
			config.IDRAC.Password,
			config.IDRAC.Port,
			config.IDRAC.UseHTTPS,
		), nil
	default:
		return nil, fmt.Errorf("unsupported BMC type: %s", config.BMCType)
	}
}
