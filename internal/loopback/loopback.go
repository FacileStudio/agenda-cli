package loopback

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"time"
)

// Listen creates a TCP listener on localhost:0 and returns the listener and port.
func Listen() (net.Listener, uint16, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, 0, err
	}

	addr := listener.Addr().(*net.TCPAddr)
	return listener, uint16(addr.Port), nil
}

// OpenBrowser attempts to open the provided URL in the default browser.
// Returns true if it succeeded, false otherwise.
func OpenBrowser(url string) bool {
	var err error
	switch {
	case isCommandAvailable("xdg-open"):
		err = runCommand("xdg-open", url)
	case isCommandAvailable("open"):
		err = runCommand("open", url)
	case isCommandAvailable("start"):
		err = runCommand("start", url)
	default:
		return false
	}
	return err == nil
}

func isCommandAvailable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func runCommand(name string, arg ...string) error {
	return exec.Command(name, arg...).Run()
}

// RandomState generates a cryptographically random state string for OAuth flows.
func RandomState() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		for i := range b {
			b[i] = byte(time.Now().UnixNano() + int64(i))
		}
	}
	return hex.EncodeToString(b)
}

// WaitForCode waits for the authorization code from the redirect URI.
func WaitForCode(listener net.Listener, expectedState string, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resultChan := make(chan waitResult, 1)
	go func() {
		resultChan <- readCode(listener)
	}()

	select {
	case result := <-resultChan:
		return validateCode(result, expectedState)
	case <-ctx.Done():
		return "", fmt.Errorf("timeout waiting for authorization code")
	}
}

type waitResult struct {
	code  string
	state string
	err   error
}

func readCode(listener net.Listener) waitResult {
	conn, err := listener.Accept()
	if err != nil {
		return waitResult{err: err}
	}
	defer conn.Close()

	req, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		return waitResult{err: err}
	}

	q := req.URL.Query()
	fmt.Fprint(conn, "HTTP/1.1 200 OK\r\nContent-Type: text/html\r\n\r\n<h1>Login successful</h1><p>You can close this window.</p>")
	return waitResult{code: q.Get("code"), state: q.Get("state")}
}

func validateCode(result waitResult, expectedState string) (string, error) {
	if result.err != nil {
		return "", result.err
	}
	if result.state != expectedState {
		return "", fmt.Errorf("invalid state: expected %s, got %s", expectedState, result.state)
	}
	return result.code, nil
}
