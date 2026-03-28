package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	cloudflareAPIBase = "https://api.cloudflare.com/client/v4"
)

// CloudflareService manages DNS records on Cloudflare
type CloudflareService struct {
	apiToken   string
	zoneID     string
	zoneName   string // e.g., "craft4free.online"
	httpClient *http.Client
}

// SRVData represents the data for an SRV record
type SRVData struct {
	Name     string `json:"name"`     // Service name (e.g., "ser4")
	Service  string `json:"service"`  // e.g., "_minecraft"
	Proto    string `json:"proto"`    // e.g., "_tcp"
	Priority int    `json:"priority"` // 0-65535
	Weight   int    `json:"weight"`   // 0-65535
	Port     int    `json:"port"`     // 0-65535
	Target   string `json:"target"`   // e.g., "worker-xxx.craft4free.online"
}

// DNSRecord represents a Cloudflare DNS record
type DNSRecord struct {
	ID      string   `json:"id,omitempty"`
	Type    string   `json:"type"`
	Name    string   `json:"name"`
	Content string   `json:"content,omitempty"`
	Data    *SRVData `json:"data,omitempty"`
	TTL     int      `json:"ttl"`
	Proxied bool     `json:"proxied,omitempty"`
}

// CloudflareResponse represents a generic Cloudflare API response
type CloudflareResponse struct {
	Success  bool              `json:"success"`
	Errors   []CloudflareError `json:"errors"`
	Messages []string          `json:"messages"`
	Result   json.RawMessage   `json:"result"`
}

// CloudflareError represents an error from Cloudflare API
type CloudflareError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// CloudflareListResponse represents a list response from Cloudflare API
type CloudflareListResponse struct {
	Success  bool              `json:"success"`
	Errors   []CloudflareError `json:"errors"`
	Messages []string          `json:"messages"`
	Result   []DNSRecord       `json:"result"`
}

// NewCloudflareService creates a new Cloudflare DNS service
func NewCloudflareService(apiToken, zoneID, zoneName string) *CloudflareService {
	return &CloudflareService{
		apiToken:   apiToken,
		zoneID:     zoneID,
		zoneName:   zoneName,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// CreateMinecraftSRVRecord creates an SRV record for a Minecraft server
// serverName: short name for the server (e.g., "ser4" or the server ID)
// workerOciID: the Oracle Cloud Instance ID for the worker
// port: the port the server is running on
func (c *CloudflareService) CreateMinecraftSRVRecord(ctx context.Context, serverName, workerOciID string, port int) error {
	// Build the target hostname: worker-{OCI_ID}.{zoneName}
	target := fmt.Sprintf("worker-%s.%s", workerOciID, c.zoneName)

	// Sanitize server name for DNS (lowercase, alphanumeric and hyphens only)
	sanitizedName := sanitizeDNSName(serverName)

	fmt.Println("Creating SRV record for", sanitizedName, "with target", target, "and port", port)
	record := DNSRecord{
		Type: "SRV",
		Name: fmt.Sprintf("_minecraft._tcp.%s", sanitizedName),
		Data: &SRVData{
			Name:     sanitizedName,
			Service:  "_minecraft",
			Proto:    "_tcp",
			Priority: 0,
			Weight:   5,
			Port:     port,
			Target:   target,
		},
		TTL: 1, // Auto TTL
	}

	body, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal DNS record: %w", err)
	}

	url := fmt.Sprintf("%s/zones/%s/dns_records", cloudflareAPIBase, c.zoneID)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var cfResp CloudflareResponse
	if err := json.Unmarshal(respBody, &cfResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if !cfResp.Success {
		errMsgs := make([]string, len(cfResp.Errors))
		for i, e := range cfResp.Errors {
			errMsgs[i] = e.Message
		}
		return fmt.Errorf("cloudflare API error: %s", strings.Join(errMsgs, ", "))
	}

	log.Printf("Created SRV record for _minecraft._tcp.%s.%s -> %s:%d", sanitizedName, c.zoneName, target, port)
	return nil
}

// DeleteMinecraftSRVRecord deletes an SRV record for a Minecraft server
func (c *CloudflareService) DeleteMinecraftSRVRecord(ctx context.Context, serverName string) error {
	sanitizedName := sanitizeDNSName(serverName)

	// First, find the record by name
	recordName := fmt.Sprintf("_minecraft._tcp.%s.%s", sanitizedName, c.zoneName)

	// List DNS records to find the one to delete
	url := fmt.Sprintf("%s/zones/%s/dns_records?type=SRV&name=%s", cloudflareAPIBase, c.zoneID, recordName)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var listResp CloudflareListResponse
	if err := json.Unmarshal(respBody, &listResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if !listResp.Success {
		errMsgs := make([]string, len(listResp.Errors))
		for i, e := range listResp.Errors {
			errMsgs[i] = e.Message
		}
		return fmt.Errorf("cloudflare API error: %s", strings.Join(errMsgs, ", "))
	}

	if len(listResp.Result) == 0 {
		log.Printf("No SRV record found for %s, nothing to delete", recordName)
		return nil
	}

	// Delete each matching record
	for _, record := range listResp.Result {
		if err := c.deleteRecord(ctx, record.ID); err != nil {
			return fmt.Errorf("failed to delete record %s: %w", record.ID, err)
		}
		log.Printf("Deleted SRV record %s (ID: %s)", recordName, record.ID)
	}

	return nil
}

// deleteRecord deletes a DNS record by ID
func (c *CloudflareService) deleteRecord(ctx context.Context, recordID string) error {
	url := fmt.Sprintf("%s/zones/%s/dns_records/%s", cloudflareAPIBase, c.zoneID, recordID)
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var cfResp CloudflareResponse
	if err := json.Unmarshal(respBody, &cfResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if !cfResp.Success {
		errMsgs := make([]string, len(cfResp.Errors))
		for i, e := range cfResp.Errors {
			errMsgs[i] = e.Message
		}
		return fmt.Errorf("cloudflare API error: %s", strings.Join(errMsgs, ", "))
	}

	return nil
}

// sanitizeDNSName converts a string to a valid DNS label
// DNS labels must be lowercase, alphanumeric, can contain hyphens (but not at start/end)
// Maximum 63 characters
func sanitizeDNSName(name string) string {
	// Convert to lowercase
	name = strings.ToLower(name)

	// Replace spaces and underscores with hyphens
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "_", "-")

	// Remove any characters that aren't alphanumeric or hyphens
	var result strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	sanitized := result.String()

	// Remove leading/trailing hyphens
	sanitized = strings.Trim(sanitized, "-")

	// Collapse multiple hyphens
	for strings.Contains(sanitized, "--") {
		sanitized = strings.ReplaceAll(sanitized, "--", "-")
	}

	// Truncate to 63 characters (DNS label limit)
	if len(sanitized) > 63 {
		sanitized = sanitized[:63]
		sanitized = strings.TrimRight(sanitized, "-")
	}

	// If empty after sanitization, use a default
	if sanitized == "" {
		sanitized = "server"
	}

	return sanitized
}

// UpdateMinecraftSRVRecord updates an existing SRV record (delete + create)
func (c *CloudflareService) UpdateMinecraftSRVRecord(ctx context.Context, serverName, workerOciID string, port int) error {
	// Delete existing record first
	if err := c.DeleteMinecraftSRVRecord(ctx, serverName); err != nil {
		log.Printf("Warning: failed to delete existing record during update: %v", err)
		// Continue anyway, the create might still work
	}

	// Create new record
	return c.CreateMinecraftSRVRecord(ctx, serverName, workerOciID, port)
}
