package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var virtualMediaCmd = &cobra.Command{
	Use:     "virtualmedia",
	Aliases: []string{"vm"},
	Short:   "Virtual media management commands",
	Long:    `Commands for managing virtual media (mounting/unmounting ISO images)`,
}

var mountCmd = &cobra.Command{
	Use:   "mount [image-url]",
	Short: "Mount virtual media",
	Long: `Mount an ISO image as virtual media. The image URL must be accessible 
from the BMC (typically an HTTP/HTTPS URL or network share).

Example:
  bmc-cli virtualmedia mount http://192.168.1.100/images/ubuntu-20.04.iso`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		imageURL := args[0]

		client, err := NewBMCClient()
		if err != nil {
			return fmt.Errorf("failed to create BMC client: %w", err)
		}

		fmt.Printf("Mounting virtual media: %s\n", imageURL)
		if err := client.MountVirtualMedia(imageURL); err != nil {
			return fmt.Errorf("failed to mount virtual media: %w", err)
		}

		fmt.Println("Virtual media mounted successfully")
		return nil
	},
}

var unmountCmd = &cobra.Command{
	Use:   "unmount",
	Short: "Unmount virtual media",
	Long:  `Unmount all virtual media from the server`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := NewBMCClient()
		if err != nil {
			return fmt.Errorf("failed to create BMC client: %w", err)
		}

		fmt.Println("Unmounting virtual media...")
		if err := client.UnmountVirtualMedia(); err != nil {
			return fmt.Errorf("failed to unmount virtual media: %w", err)
		}

		fmt.Println("Virtual media unmounted successfully")
		return nil
	},
}

var listMediaCmd = &cobra.Command{
	Use:   "list",
	Short: "List virtual media status",
	Long:  `Lists all virtual media slots and their current status`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := NewBMCClient()
		if err != nil {
			return fmt.Errorf("failed to create BMC client: %w", err)
		}

		fmt.Println("Retrieving virtual media information...")
		vmList, err := client.GetVirtualMedia()
		if err != nil {
			return fmt.Errorf("failed to get virtual media info: %w", err)
		}

		if len(vmList) == 0 {
			fmt.Println("No virtual media slots found")
			return nil
		}

		fmt.Printf("%-15s %-15s %-10s %-10s %s\n", "Name", "Media Types", "Connected", "Inserted", "Image")
		fmt.Println("---------------------------------------------------------------------------------")

		for _, vm := range vmList {
			mediaTypes := "None"
			if len(vm.MediaTypes) > 0 {
				mediaTypes = vm.MediaTypes[0]
				for _, mt := range vm.MediaTypes[1:] {
					mediaTypes += ", " + mt
				}
			}

			connected := "No"
			if vm.Connected {
				connected = "Yes"
			}

			inserted := "No"
			if vm.Inserted {
				inserted = "Yes"
			}

			image := vm.Image
			if image == "" {
				image = "-"
			}

			fmt.Printf("%-15s %-15s %-10s %-10s %s\n",
				vm.Name, mediaTypes, connected, inserted, image)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(virtualMediaCmd)
	virtualMediaCmd.AddCommand(mountCmd)
	virtualMediaCmd.AddCommand(unmountCmd)
	virtualMediaCmd.AddCommand(listMediaCmd)
}
