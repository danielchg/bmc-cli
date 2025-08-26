package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// IDRACClient represents an iDRAC API client
type IDRACClient struct {
	baseURL    string
	username   string
	password   string
	httpClient *http.Client
}

// NewIDRACClient creates a new iDRAC client
func NewIDRACClient(host, username, password string, port int, useHTTPS bool) BMCClient {
	scheme := "http"
	if useHTTPS {
		scheme = "https"
	}

	baseURL := fmt.Sprintf("%s://%s:%d", scheme, host, port)

	// Create HTTP client with custom transport for SSL/TLS
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // For self-signed certificates
		},
	}

	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}

	return &IDRACClient{
		baseURL:    baseURL,
		username:   username,
		password:   password,
		httpClient: client,
	}
}

// makeRequest makes an authenticated HTTP request to the iDRAC API
func (c *IDRACClient) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	url := c.baseURL + endpoint
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if verbose {
		fmt.Printf("Making %s request to %s\n", method, url)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	return resp, nil
}

// GetSystemInfo retrieves basic system information
func (c *IDRACClient) GetSystemInfo() (*SystemInfo, error) {
	resp, err := c.makeRequest("GET", "/redfish/v1/Systems/System.Embedded.1", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var systemInfo SystemInfo
	if err := json.NewDecoder(resp.Body).Decode(&systemInfo); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &systemInfo, nil
}

// SetPowerState changes the server power state
func (c *IDRACClient) SetPowerState(state PowerState) error {
	powerRequest := PowerRequest{
		ResetType: string(state),
	}

	resp, err := c.makeRequest("POST", "/redfish/v1/Systems/System.Embedded.1/Actions/ComputerSystem.Reset", powerRequest)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusAccepted {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("power operation failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// GetVirtualMedia lists available virtual media slots
func (c *IDRACClient) GetVirtualMedia() ([]VirtualMediaInfo, error) {
	resp, err := c.makeRequest("GET", "/redfish/v1/Managers/iDRAC.Embedded.1/VirtualMedia", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var result struct {
		Members []struct {
			OdataID string `json:"@odata.id"`
		} `json:"Members"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	var virtualMediaList []VirtualMediaInfo
	for _, member := range result.Members {
		vmInfo, err := c.getVirtualMediaInfo(member.OdataID)
		if err != nil {
			continue // Skip this media slot if we can't get info
		}
		virtualMediaList = append(virtualMediaList, *vmInfo)
	}

	return virtualMediaList, nil
}

// getVirtualMediaInfo gets information about a specific virtual media slot
func (c *IDRACClient) getVirtualMediaInfo(odataID string) (*VirtualMediaInfo, error) {
	resp, err := c.makeRequest("GET", odataID, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var vmInfo VirtualMediaInfo
	if err := json.NewDecoder(resp.Body).Decode(&vmInfo); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &vmInfo, nil
}

// MountVirtualMedia mounts an image to the first available CD/DVD virtual media slot
func (c *IDRACClient) MountVirtualMedia(imageURL string) error {
	// Get available virtual media slots
	vmList, err := c.GetVirtualMedia()
	if err != nil {
		return fmt.Errorf("error getting virtual media info: %w", err)
	}

	// Find the first CD/DVD slot
	var targetSlot string
	for _, vm := range vmList {
		for _, mediaType := range vm.MediaTypes {
			if mediaType == "CD" || mediaType == "DVD" {
				// iDRAC uses different naming convention
				if strings.Contains(vm.Name, "CD") || strings.Contains(vm.Name, "DVD") {
					targetSlot = fmt.Sprintf("/redfish/v1/Managers/iDRAC.Embedded.1/VirtualMedia/%s", vm.Name)
					break
				}
			}
		}
		if targetSlot != "" {
			break
		}
	}

	if targetSlot == "" {
		return fmt.Errorf("no CD/DVD virtual media slot found")
	}

	// For iDRAC, we need to first insert the media, then connect
	// First, eject any existing media
	ejectRequest := map[string]interface{}{
		"Inserted": false,
	}

	resp, err := c.makeRequest("PATCH", targetSlot, ejectRequest)
	if err == nil {
		resp.Body.Close()
	}

	// Now mount the new image
	mountRequest := map[string]interface{}{
		"Image":    imageURL,
		"Inserted": true,
	}

	resp, err = c.makeRequest("PATCH", targetSlot, mountRequest)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusAccepted {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("virtual media mount failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// UnmountVirtualMedia unmounts virtual media from all slots
func (c *IDRACClient) UnmountVirtualMedia() error {
	vmList, err := c.GetVirtualMedia()
	if err != nil {
		return fmt.Errorf("error getting virtual media info: %w", err)
	}

	for _, vm := range vmList {
		if vm.Inserted {
			unmountRequest := map[string]interface{}{
				"Inserted": false,
			}

			endpoint := fmt.Sprintf("/redfish/v1/Managers/iDRAC.Embedded.1/VirtualMedia/%s", vm.Name)
			resp, err := c.makeRequest("PATCH", endpoint, unmountRequest)
			if err != nil {
				continue // Continue with other slots
			}
			resp.Body.Close()
		}
	}

	return nil
}
