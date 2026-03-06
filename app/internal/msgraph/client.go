// Package msgraph provides a lightweight Microsoft Graph API client.
// Uses raw HTTP with OAuth2 client credentials — no heavy SDK dependency.
// Supports commercial, GCC High, and DoD national cloud endpoints.
package msgraph

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// CloudConfig holds endpoints for a specific Microsoft cloud.
type CloudConfig struct {
	LoginBase string // e.g. "https://login.microsoftonline.com"
	GraphBase string // e.g. "https://graph.microsoft.com"
}

var clouds = map[string]CloudConfig{
	"commercial": {
		LoginBase: "https://login.microsoftonline.com",
		GraphBase: "https://graph.microsoft.com",
	},
	"gcc-high": {
		LoginBase: "https://login.microsoftonline.us",
		GraphBase: "https://graph.microsoft.us",
	},
	"dod": {
		LoginBase: "https://login.microsoftonline.us",
		GraphBase: "https://dod-graph.microsoft.us",
	},
}

// Client is a Microsoft Graph API client using client credentials OAuth2 flow.
type Client struct {
	tenantID     string
	clientID     string
	clientSecret string
	cloud        CloudConfig
	httpClient   *http.Client

	mu    sync.RWMutex
	token string
	expAt time.Time
}

// NewClient creates a Graph client for the specified cloud environment.
// cloud: "commercial", "gcc-high", or "dod"
func NewClient(tenantID, clientID, clientSecret, cloud string) *Client {
	cfg, ok := clouds[cloud]
	if !ok {
		cfg = clouds["commercial"]
	}

	return &Client{
		tenantID:     tenantID,
		clientID:     clientID,
		clientSecret: clientSecret,
		cloud:        cfg,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
	}
}

// IsConfigured returns true if credentials are set.
func (c *Client) IsConfigured() bool {
	return c.tenantID != "" && c.clientID != "" && c.clientSecret != ""
}

// getToken returns a valid access token, refreshing if needed.
func (c *Client) getToken() (string, error) {
	c.mu.RLock()
	if c.token != "" && time.Now().Before(c.expAt) {
		defer c.mu.RUnlock()
		return c.token, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after acquiring write lock
	if c.token != "" && time.Now().Before(c.expAt) {
		return c.token, nil
	}

	tokenURL := fmt.Sprintf("%s/%s/oauth2/v2.0/token", c.cloud.LoginBase, c.tenantID)
	scope := c.cloud.GraphBase + "/.default"

	data := url.Values{
		"client_id":     {c.clientID},
		"client_secret": {c.clientSecret},
		"scope":         {scope},
		"grant_type":    {"client_credentials"},
	}

	resp, err := c.httpClient.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("token request returned %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("parse token response: %w", err)
	}

	c.token = tokenResp.AccessToken
	// Refresh 5 minutes before expiry
	c.expAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn-300) * time.Second)

	slog.Info("Microsoft Graph token acquired", "expiresIn", tokenResp.ExpiresIn)
	return c.token, nil
}

// Get performs an authenticated GET request to the Graph API.
// path should start with "/" (e.g. "/users/user@example.com/calendarView")
func (c *Client) Get(path string, params map[string]string) ([]byte, error) {
	token, err := c.getToken()
	if err != nil {
		return nil, err
	}

	u := c.cloud.GraphBase + "/v1.0" + path
	if len(params) > 0 {
		q := url.Values{}
		for k, v := range params {
			q.Set(k, v)
		}
		u += "?" + q.Encode()
	}

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("graph request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("graph %s returned %d: %s", path, resp.StatusCode, truncate(string(body), 200))
	}

	return body, nil
}

// Post performs an authenticated POST request to the Graph API.
func (c *Client) Post(path string, payload interface{}) ([]byte, error) {
	token, err := c.getToken()
	if err != nil {
		return nil, err
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	u := c.cloud.GraphBase + "/v1.0" + path
	req, err := http.NewRequest("POST", u, strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("graph POST failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("graph POST %s returned %d: %s", path, resp.StatusCode, truncate(string(body), 200))
	}

	return body, nil
}

// TestConnection verifies the Graph API credentials work.
func (c *Client) TestConnection() error {
	_, err := c.getToken()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Verify token works by calling a minimal endpoint
	_, err = c.Get("/organization", map[string]string{"$select": "id"})
	if err != nil {
		return fmt.Errorf("graph API test failed: %w", err)
	}
	return nil
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
