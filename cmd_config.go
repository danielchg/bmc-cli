package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration management commands",
	Long:  `Commands for managing configuration files and settings`,
}

var generateConfigCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a sample configuration file",
	Long: `Generate a sample YAML configuration file with all available options.
This creates a config.yaml file in the current directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := createSampleConfig(); err != nil {
			return fmt.Errorf("failed to generate config: %w", err)
		}
		return nil
	},
}

var showConfigCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Display the current configuration values (without showing passwords)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Current Configuration:")
		fmt.Println("=====================")
		fmt.Printf("Host: %s\n", config.ILO.Host)
		fmt.Printf("Username: %s\n", config.ILO.Username)
		fmt.Printf("Port: %d\n", config.ILO.Port)
		fmt.Printf("Use HTTPS: %t\n", config.ILO.UseHTTPS)

		if config.ILO.Password != "" {
			fmt.Printf("Password: %s\n", "***configured***")
		} else {
			fmt.Printf("Password: %s\n", "***not configured***")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(generateConfigCmd)
	configCmd.AddCommand(showConfigCmd)
}
