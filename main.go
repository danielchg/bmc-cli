package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "bmc-cli",
	Short: "A CLI tool to manage BMC operations",
	Long: `bmc-cli is a command-line tool for managing Baseboard Management Controllers (BMCs)
including HP iLO (Integrated Lights-Out) and DELL iDRAC systems. It provides functionality to:
- Power on/off servers
- Mount virtual media
- Manage server configurations

Supports both HPE iLO and DELL iDRAC BMCs via Redfish API.
Configuration can be provided via YAML file or environment variables.`,
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

func initConfig() {
	if err := loadConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	if verbose {
		log.SetOutput(os.Stdout)
		switch config.BMCType {
		case BMCTypeILO:
			fmt.Printf("Connected to iLO at %s:%d\n", config.ILO.Host, config.ILO.Port)
		case BMCTypeIDRAC:
			fmt.Printf("Connected to iDRAC at %s:%d\n", config.IDRAC.Host, config.IDRAC.Port)
		}
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
