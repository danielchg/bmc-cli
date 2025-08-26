package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIDRACClient_GetSystemInfo(t *testing.T) {
	// Mock response
	mockResponse := SystemInfo{
		PowerState: "On",
		Status: struct {
			Health string `json:"Health"`
			State  string `json:"State"`
		}{
			Health: "OK",
			State:  "Enabled",
		},
	}

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/redfish/v1/Systems/System.Embedded.1" {
			t.Errorf("Expected path '/redfish/v1/Systems/System.Embedded.1', got: %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got: %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Create client with test server URL
	client := &IDRACClient{
		baseURL:    server.URL,
		username:   "root",
		password:   "calvin",
		httpClient: server.Client(),
	}

	// Test GetSystemInfo
	systemInfo, err := client.GetSystemInfo()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if systemInfo.PowerState != "On" {
		t.Errorf("Expected PowerState 'On', got: %s", systemInfo.PowerState)
	}
	if systemInfo.Status.Health != "OK" {
		t.Errorf("Expected Health 'OK', got: %s", systemInfo.Status.Health)
	}
}

func TestIDRACClient_SetPowerState_On(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset" {
			t.Errorf("Expected path '/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset', got: %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got: %s", r.Method)
		}

		// Verify request body
		var powerRequest PowerRequest
		if err := json.NewDecoder(r.Body).Decode(&powerRequest); err != nil {
			t.Errorf("Error decoding request body: %v", err)
		}
		if powerRequest.ResetType != "On" {
			t.Errorf("Expected ResetType 'On', got: %s", powerRequest.ResetType)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create client with test server URL
	client := &IDRACClient{
		baseURL:    server.URL,
		username:   "root",
		password:   "calvin",
		httpClient: server.Client(),
	}

	// Test SetPowerState
	err := client.SetPowerState(PowerStateOn)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestIDRACClient_SetPowerState_Off(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset" {
			t.Errorf("Expected path '/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset', got: %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got: %s", r.Method)
		}

		// Verify request body
		var powerRequest PowerRequest
		if err := json.NewDecoder(r.Body).Decode(&powerRequest); err != nil {
			t.Errorf("Error decoding request body: %v", err)
		}
		if powerRequest.ResetType != "ForceOff" {
			t.Errorf("Expected ResetType 'ForceOff', got: %s", powerRequest.ResetType)
		}

		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	// Create client with test server URL
	client := &IDRACClient{
		baseURL:    server.URL,
		username:   "root",
		password:   "calvin",
		httpClient: server.Client(),
	}

	// Test SetPowerState
	err := client.SetPowerState(PowerStateOff)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestIDRACClient_GetVirtualMedia(t *testing.T) {
	// Mock responses
	membersResponse := struct {
		Members []struct {
			OdataID string `json:"@odata.id"`
		} `json:"Members"`
	}{
		Members: []struct {
			OdataID string `json:"@odata.id"`
		}{
			{OdataID: "/redfish/v1/Managers/iDRAC.Embedded.1/VirtualMedia/CD"},
			{OdataID: "/redfish/v1/Managers/iDRAC.Embedded.1/VirtualMedia/RemovableDisk"},
		},
	}

	vmInfoCD := VirtualMediaInfo{
		Name:       "CD",
		MediaTypes: []string{"CD", "DVD"},
		Connected:  false,
		Inserted:   false,
		Image:      "",
	}

	vmInfoDisk := VirtualMediaInfo{
		Name:       "RemovableDisk",
		MediaTypes: []string{"USBStick"},
		Connected:  false,
		Inserted:   false,
		Image:      "",
	}

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Path == "/redfish/v1/Managers/iDRAC.Embedded.1/VirtualMedia" {
			json.NewEncoder(w).Encode(membersResponse)
		} else if r.URL.Path == "/redfish/v1/Managers/iDRAC.Embedded.1/VirtualMedia/CD" {
			json.NewEncoder(w).Encode(vmInfoCD)
		} else if r.URL.Path == "/redfish/v1/Managers/iDRAC.Embedded.1/VirtualMedia/RemovableDisk" {
			json.NewEncoder(w).Encode(vmInfoDisk)
		} else {
			t.Errorf("Unexpected path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Create client with test server URL
	client := &IDRACClient{
		baseURL:    server.URL,
		username:   "root",
		password:   "calvin",
		httpClient: server.Client(),
	}

	// Test GetVirtualMedia
	vmList, err := client.GetVirtualMedia()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(vmList) != 2 {
		t.Errorf("Expected 2 virtual media items, got: %d", len(vmList))
	}
	if vmList[0].Name != "CD" {
		t.Errorf("Expected name 'CD', got: %s", vmList[0].Name)
	}
}

func TestIDRACClient_MountVirtualMedia(t *testing.T) {
	// Mock responses
	membersResponse := struct {
		Members []struct {
			OdataID string `json:"@odata.id"`
		} `json:"Members"`
	}{
		Members: []struct {
			OdataID string `json:"@odata.id"`
		}{
			{OdataID: "/redfish/v1/Managers/iDRAC.Embedded.1/VirtualMedia/CD"},
		},
	}

	vmInfo := VirtualMediaInfo{
		Name:       "CD",
		MediaTypes: []string{"CD", "DVD"},
		Connected:  false,
		Inserted:   false,
		Image:      "",
	}

	requestCount := 0

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Path == "/redfish/v1/Managers/iDRAC.Embedded.1/VirtualMedia" {
			json.NewEncoder(w).Encode(membersResponse)
		} else if r.URL.Path == "/redfish/v1/Managers/iDRAC.Embedded.1/VirtualMedia/CD" {
			if r.Method == "GET" {
				json.NewEncoder(w).Encode(vmInfo)
			} else if r.Method == "PATCH" {
				requestCount++

				// Verify mount request
				var mountRequest map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&mountRequest); err != nil {
					t.Errorf("Error decoding request body: %v", err)
				}

				if requestCount == 2 { // Second PATCH request (actual mount)
					if mountRequest["Image"] != "http://example.com/image.iso" {
						t.Errorf("Expected image URL 'http://example.com/image.iso', got: %v", mountRequest["Image"])
					}
					if mountRequest["Inserted"] != true {
						t.Error("Expected Inserted to be true")
					}
				}

				w.WriteHeader(http.StatusOK)
			}
		} else {
			t.Errorf("Unexpected path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Create client with test server URL
	client := &IDRACClient{
		baseURL:    server.URL,
		username:   "root",
		password:   "calvin",
		httpClient: server.Client(),
	}

	// Test MountVirtualMedia
	err := client.MountVirtualMedia("http://example.com/image.iso")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestIDRACClient_UnmountVirtualMedia(t *testing.T) {
	// Mock responses
	membersResponse := struct {
		Members []struct {
			OdataID string `json:"@odata.id"`
		} `json:"Members"`
	}{
		Members: []struct {
			OdataID string `json:"@odata.id"`
		}{
			{OdataID: "/redfish/v1/Managers/iDRAC.Embedded.1/VirtualMedia/CD"},
		},
	}

	vmInfo := VirtualMediaInfo{
		Name:       "CD",
		MediaTypes: []string{"CD", "DVD"},
		Connected:  true,
		Inserted:   true,
		Image:      "http://example.com/image.iso",
	}

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Path == "/redfish/v1/Managers/iDRAC.Embedded.1/VirtualMedia" {
			json.NewEncoder(w).Encode(membersResponse)
		} else if r.URL.Path == "/redfish/v1/Managers/iDRAC.Embedded.1/VirtualMedia/CD" {
			if r.Method == "GET" {
				json.NewEncoder(w).Encode(vmInfo)
			} else if r.Method == "PATCH" {
				// Verify unmount request
				var unmountRequest map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&unmountRequest); err != nil {
					t.Errorf("Error decoding request body: %v", err)
				}
				if unmountRequest["Inserted"] != false {
					t.Error("Expected Inserted to be false")
				}
				w.WriteHeader(http.StatusOK)
			}
		} else {
			t.Errorf("Unexpected path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Create client with test server URL
	client := &IDRACClient{
		baseURL:    server.URL,
		username:   "root",
		password:   "calvin",
		httpClient: server.Client(),
	}

	// Test UnmountVirtualMedia
	err := client.UnmountVirtualMedia()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestIDRACClient_ErrorHandling(t *testing.T) {
	// Create test server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
	}))
	defer server.Close()

	// Create client with test server URL
	client := &IDRACClient{
		baseURL:    server.URL,
		username:   "root",
		password:   "calvin",
		httpClient: server.Client(),
	}

	// Test error handling for GetSystemInfo
	_, err := client.GetSystemInfo()
	if err == nil {
		t.Error("Expected error for 401 response, got nil")
	}

	// Test error handling for SetPowerState
	err = client.SetPowerState(PowerStateOn)
	if err == nil {
		t.Error("Expected error for 401 response, got nil")
	}
}

func TestIDRACClient_NoVirtualMediaFound(t *testing.T) {
	// Mock responses with no CD/DVD slots
	membersResponse := struct {
		Members []struct {
			OdataID string `json:"@odata.id"`
		} `json:"Members"`
	}{
		Members: []struct {
			OdataID string `json:"@odata.id"`
		}{
			{OdataID: "/redfish/v1/Managers/iDRAC.Embedded.1/VirtualMedia/RemovableDisk"},
		},
	}

	vmInfo := VirtualMediaInfo{
		Name:       "RemovableDisk",
		MediaTypes: []string{"USBStick"},
		Connected:  false,
		Inserted:   false,
		Image:      "",
	}

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Path == "/redfish/v1/Managers/iDRAC.Embedded.1/VirtualMedia" {
			json.NewEncoder(w).Encode(membersResponse)
		} else if strings.HasPrefix(r.URL.Path, "/redfish/v1/Managers/iDRAC.Embedded.1/VirtualMedia/") {
			json.NewEncoder(w).Encode(vmInfo)
		}
	}))
	defer server.Close()

	// Create client with test server URL
	client := &IDRACClient{
		baseURL:    server.URL,
		username:   "root",
		password:   "calvin",
		httpClient: server.Client(),
	}

	// Test MountVirtualMedia when no CD/DVD slot is available
	err := client.MountVirtualMedia("http://example.com/image.iso")
	if err == nil {
		t.Error("Expected error when no CD/DVD slot found, got nil")
	}
	if !strings.Contains(err.Error(), "no CD/DVD virtual media slot found") {
		t.Errorf("Expected error about no CD/DVD slot, got: %v", err)
	}
}
