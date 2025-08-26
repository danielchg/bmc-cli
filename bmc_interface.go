package main

// BMCClient interface defines common operations for all BMC types
type BMCClient interface {
	GetSystemInfo() (*SystemInfo, error)
	SetPowerState(state PowerState) error
	GetVirtualMedia() ([]VirtualMediaInfo, error)
	MountVirtualMedia(imageURL string) error
	UnmountVirtualMedia() error
}

// BMCType represents the type of BMC hardware
type BMCType string

const (
	BMCTypeILO   BMCType = "ilo"
	BMCTypeIDRAC BMCType = "idrac"
)
