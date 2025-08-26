package main

import (
	"testing"
)

func TestBMCInterface_ILOClient(t *testing.T) {
	// Test that ILOClient implements BMCClient interface
	var client BMCClient
	iloClient := &ILOClient{
		baseURL:  "https://test.example.com",
		username: "admin",
		password: "password",
	}

	// This should compile without errors if ILOClient implements BMCClient
	client = iloClient
	if client == nil {
		t.Error("ILOClient should implement BMCClient interface")
	}
}

func TestBMCInterface_IDRACClient(t *testing.T) {
	// Test that IDRACClient implements BMCClient interface
	var client BMCClient
	idracClient := &IDRACClient{
		baseURL:  "https://test.example.com",
		username: "root",
		password: "calvin",
	}

	// This should compile without errors if IDRACClient implements BMCClient
	client = idracClient
	if client == nil {
		t.Error("IDRACClient should implement BMCClient interface")
	}
}

func TestBMCTypes(t *testing.T) {
	// Test BMC type constants
	if BMCTypeILO != "ilo" {
		t.Errorf("Expected BMCTypeILO to be 'ilo', got: %s", BMCTypeILO)
	}
	if BMCTypeIDRAC != "idrac" {
		t.Errorf("Expected BMCTypeIDRAC to be 'idrac', got: %s", BMCTypeIDRAC)
	}
}

func TestPowerStates(t *testing.T) {
	// Test power state constants
	if PowerStateOn != "On" {
		t.Errorf("Expected PowerStateOn to be 'On', got: %s", PowerStateOn)
	}
	if PowerStateOff != "ForceOff" {
		t.Errorf("Expected PowerStateOff to be 'ForceOff', got: %s", PowerStateOff)
	}
}
