package main

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestValidateConfig_ILO(t *testing.T) {
	// Test valid iLO configuration
	config = Config{
		BMCType: BMCTypeILO,
		ILO: ILOConfig{
			Host:     "192.168.1.100",
			Username: "admin",
			Password: "password",
			Port:     443,
			UseHTTPS: true,
		},
	}

	err := validateConfig()
	if err != nil {
		t.Errorf("Expected no error for valid iLO config, got: %v", err)
	}
}

func TestValidateConfig_IDRAC(t *testing.T) {
	// Test valid iDRAC configuration
	config = Config{
		BMCType: BMCTypeIDRAC,
		IDRAC: IDRACConfig{
			Host:     "192.168.1.101",
			Username: "root",
			Password: "calvin",
			Port:     443,
			UseHTTPS: true,
		},
	}

	err := validateConfig()
	if err != nil {
		t.Errorf("Expected no error for valid iDRAC config, got: %v", err)
	}
}

func TestValidateConfig_MissingILOHost(t *testing.T) {
	config = Config{
		BMCType: BMCTypeILO,
		ILO: ILOConfig{
			Host:     "", // Missing host
			Username: "admin",
			Password: "password",
		},
	}

	err := validateConfig()
	if err == nil {
		t.Error("Expected error for missing iLO host, got nil")
	}
}

func TestValidateConfig_MissingIDRACHost(t *testing.T) {
	config = Config{
		BMCType: BMCTypeIDRAC,
		IDRAC: IDRACConfig{
			Host:     "", // Missing host
			Username: "root",
			Password: "calvin",
		},
	}

	err := validateConfig()
	if err == nil {
		t.Error("Expected error for missing iDRAC host, got nil")
	}
}

func TestValidateConfig_UnsupportedBMCType(t *testing.T) {
	config = Config{
		BMCType: "unsupported",
	}

	err := validateConfig()
	if err == nil {
		t.Error("Expected error for unsupported BMC type, got nil")
	}
}

func TestNewBMCClient_ILO(t *testing.T) {
	config = Config{
		BMCType: BMCTypeILO,
		ILO: ILOConfig{
			Host:     "192.168.1.100",
			Username: "admin",
			Password: "password",
			Port:     443,
			UseHTTPS: true,
		},
	}

	client, err := NewBMCClient()
	if err != nil {
		t.Errorf("Expected no error creating iLO client, got: %v", err)
	}
	if client == nil {
		t.Error("Expected client to be created, got nil")
	}
}

func TestNewBMCClient_IDRAC(t *testing.T) {
	config = Config{
		BMCType: BMCTypeIDRAC,
		IDRAC: IDRACConfig{
			Host:     "192.168.1.101",
			Username: "root",
			Password: "calvin",
			Port:     443,
			UseHTTPS: true,
		},
	}

	client, err := NewBMCClient()
	if err != nil {
		t.Errorf("Expected no error creating iDRAC client, got: %v", err)
	}
	if client == nil {
		t.Error("Expected client to be created, got nil")
	}
}

func TestNewBMCClient_UnsupportedType(t *testing.T) {
	config = Config{
		BMCType: "unsupported",
	}

	client, err := NewBMCClient()
	if err == nil {
		t.Error("Expected error for unsupported BMC type, got nil")
	}
	if client != nil {
		t.Error("Expected nil client for unsupported type, got client")
	}
}

func TestLoadConfigFromEnvironment_ILO(t *testing.T) {
	// Save original config
	originalConfig := config

	// Set environment variables
	os.Setenv("BMC_TYPE", "ilo")
	os.Setenv("ILO_HOST", "192.168.1.100")
	os.Setenv("ILO_USERNAME", "admin")
	os.Setenv("ILO_PASSWORD", "password")
	os.Setenv("ILO_PORT", "443")
	os.Setenv("ILO_USE_HTTPS", "true")

	// Clear viper
	viper.Reset()

	// Load config
	err := loadConfig()
	if err != nil {
		t.Errorf("Expected no error loading config from environment, got: %v", err)
	}

	// Verify config values
	if config.BMCType != BMCTypeILO {
		t.Errorf("Expected BMC type 'ilo', got: %v", config.BMCType)
	}
	if config.ILO.Host != "192.168.1.100" {
		t.Errorf("Expected host '192.168.1.100', got: %v", config.ILO.Host)
	}
	if config.ILO.Username != "admin" {
		t.Errorf("Expected username 'admin', got: %v", config.ILO.Username)
	}

	// Cleanup
	os.Unsetenv("BMC_TYPE")
	os.Unsetenv("ILO_HOST")
	os.Unsetenv("ILO_USERNAME")
	os.Unsetenv("ILO_PASSWORD")
	os.Unsetenv("ILO_PORT")
	os.Unsetenv("ILO_USE_HTTPS")
	config = originalConfig
}

func TestLoadConfigFromEnvironment_IDRAC(t *testing.T) {
	// Save original config
	originalConfig := config

	// Set environment variables
	os.Setenv("BMC_TYPE", "idrac")
	os.Setenv("IDRAC_HOST", "192.168.1.101")
	os.Setenv("IDRAC_USERNAME", "root")
	os.Setenv("IDRAC_PASSWORD", "calvin")
	os.Setenv("IDRAC_PORT", "443")
	os.Setenv("IDRAC_USE_HTTPS", "true")

	// Clear viper
	viper.Reset()

	// Load config
	err := loadConfig()
	if err != nil {
		t.Errorf("Expected no error loading config from environment, got: %v", err)
	}

	// Verify config values
	if config.BMCType != BMCTypeIDRAC {
		t.Errorf("Expected BMC type 'idrac', got: %v", config.BMCType)
	}
	if config.IDRAC.Host != "192.168.1.101" {
		t.Errorf("Expected host '192.168.1.101', got: %v", config.IDRAC.Host)
	}
	if config.IDRAC.Username != "root" {
		t.Errorf("Expected username 'root', got: %v", config.IDRAC.Username)
	}

	// Cleanup
	os.Unsetenv("BMC_TYPE")
	os.Unsetenv("IDRAC_HOST")
	os.Unsetenv("IDRAC_USERNAME")
	os.Unsetenv("IDRAC_PASSWORD")
	os.Unsetenv("IDRAC_PORT")
	os.Unsetenv("IDRAC_USE_HTTPS")
	config = originalConfig
}
