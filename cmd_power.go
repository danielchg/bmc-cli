package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var powerCmd = &cobra.Command{
	Use:   "power",
	Short: "Power management commands",
	Long:  `Commands for managing server power state (on, off, status)`,
}

var powerOnCmd = &cobra.Command{
	Use:   "on",
	Short: "Power on the server",
	Long:  `Powers on the server via iLO BMC`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := NewILOClient(
			config.ILO.Host,
			config.ILO.Username,
			config.ILO.Password,
			config.ILO.Port,
			config.ILO.UseHTTPS,
		)

		fmt.Println("Powering on server...")
		if err := client.SetPowerState(PowerStateOn); err != nil {
			return fmt.Errorf("failed to power on server: %w", err)
		}

		fmt.Println("Server power on command sent successfully")
		return nil
	},
}

var powerOffCmd = &cobra.Command{
	Use:   "off",
	Short: "Power off the server",
	Long:  `Forces the server to power off via iLO BMC`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := NewILOClient(
			config.ILO.Host,
			config.ILO.Username,
			config.ILO.Password,
			config.ILO.Port,
			config.ILO.UseHTTPS,
		)

		fmt.Println("Powering off server...")
		if err := client.SetPowerState(PowerStateOff); err != nil {
			return fmt.Errorf("failed to power off server: %w", err)
		}

		fmt.Println("Server power off command sent successfully")
		return nil
	},
}

var powerStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check server power status",
	Long:  `Retrieves the current power state and health status of the server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := NewILOClient(
			config.ILO.Host,
			config.ILO.Username,
			config.ILO.Password,
			config.ILO.Port,
			config.ILO.UseHTTPS,
		)

		fmt.Println("Checking server status...")
		systemInfo, err := client.GetSystemInfo()
		if err != nil {
			return fmt.Errorf("failed to get system info: %w", err)
		}

		fmt.Printf("Power State: %s\n", systemInfo.PowerState)
		fmt.Printf("Health: %s\n", systemInfo.Status.Health)
		fmt.Printf("State: %s\n", systemInfo.Status.State)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(powerCmd)
	powerCmd.AddCommand(powerOnCmd)
	powerCmd.AddCommand(powerOffCmd)
	powerCmd.AddCommand(powerStatusCmd)
}
