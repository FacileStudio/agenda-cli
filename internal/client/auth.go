package client

// AuthConfigResponse represents the server's auth configuration.
type AuthConfigResponse struct {
	SSOOnly     bool `json:"sso_only"`
	OIDCEnabled bool `json:"oidc_enabled"`
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
