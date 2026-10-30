package client

import (
	"fmt"
	"io"
	"net/http"
)

// ExpectOK closes the response and turns an error status into an error that
// carries the server's message, for calls whose body is of no interest.
func (c *Client) ExpectOK(resp *http.Response) error {
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		data, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(data))
	}
	return nil
}
