package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/FacileStudio/agenda-cli/internal/ui"
)

type Client struct {
	config     *Config
	httpClient *http.Client
}

func NewClient(config *Config) *Client {
	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) Get(path string) (*http.Response, error) {
	return c.doRequest("GET", path, nil)
}

// GetWithQuery performs a GET request with query parameters.
func (c *Client) GetWithQuery(path string, params map[string]string) (*http.Response, error) {
	if err := c.config.Validate(); err != nil {
		return nil, err
	}

	u, err := url.Parse(c.config.ServerURL + path)
	if err != nil {
		return nil, fmt.Errorf("parse URL: %w", err)
	}

	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(context.Background(), "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	if c.config.Token != "" {
		req.Header.Set("Authorization", c.config.AuthHeader())
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

func (c *Client) Post(path string, body interface{}) (*http.Response, error) {
	return c.doRequest("POST", path, body)
}

// AuthConfig fetches the server's authentication configuration.
func (c *Client) AuthConfig() (*AuthConfigResponse, error) {
	resp, err := c.Get("/auth/config")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var config AuthConfigResponse
	if err := c.DecodeJSON(resp, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// Exchange exchanges a device code for an access token.
func (c *Client) Exchange(deviceCode string) (*string, error) {
	resp, err := c.Post("/auth/device/exchange", map[string]string{
		"device_code": deviceCode,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var respStruct struct {
		Token string `json:"token"`
	}
	if err := c.DecodeJSON(resp, &respStruct); err != nil {
		return nil, err
	}
	return &respStruct.Token, nil
}

// AuthConfigResponse represents the server's auth configuration.
type AuthConfigResponse struct {
	SSOOnly    bool `json:"sso_only"`
	OIDCEnabled bool `json:"oidc_enabled"`
}

func (c *Client) Put(path string, body interface{}) (*http.Response, error) {
	return c.doRequest("PUT", path, body)
}

func (c *Client) Delete(path string) (*http.Response, error) {
	return c.doRequest("DELETE", path, nil)
}

func (c *Client) doRequest(method, path string, body interface{}) (*http.Response, error) {
	if err := c.config.Validate(); err != nil {
		return nil, err
	}

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(context.Background(), method, c.config.ServerURL+path, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	if c.config.Token != "" {
		req.Header.Set("Authorization", c.config.AuthHeader())
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

func (c *Client) DecodeJSON(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		data, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(data))
	}

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return err
	}

	return nil
}

func (c *Client) PrintJSON(v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		ui.PrintError(err)
		return
	}
	fmt.Fprintln(os.Stdout, string(data))
}
