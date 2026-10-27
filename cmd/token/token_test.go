package token

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setupTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(handler)
	os.Setenv("AGENDA_SERVER_URL", server.URL)
	os.Setenv("AGENDA_TOKEN", "test-token")
	t.Cleanup(func() {
		server.Close()
		os.Unsetenv("AGENDA_SERVER_URL")
		os.Unsetenv("AGENDA_TOKEN")
	})
	return server
}

func TestTokenStatusHasToken(t *testing.T) {
	setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/me/api-token" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"has_token":  true,
			"name":       "cli-token",
			"created_at": "2026-10-27T00:00:00Z",
		})
	})

	if err := runStatus(true); err != nil {
		t.Fatalf("runStatus failed: %v", err)
	}
	if err := runStatus(false); err != nil {
		t.Fatalf("runStatus failed: %v", err)
	}
}

func TestTokenStatusNoToken(t *testing.T) {
	setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"has_token": false})
	})

	if err := runStatus(true); err != nil {
		t.Fatalf("runStatus failed: %v", err)
	}
	if err := runStatus(false); err != nil {
		t.Fatalf("runStatus failed: %v", err)
	}
}

func TestTokenCreate(t *testing.T) {
	setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"token":      "tok_123456",
			"name":       "my-token",
			"created_at": "2026-10-27T00:00:00Z",
		})
	})

	if err := runCreate("my-token", true); err != nil {
		t.Fatalf("runCreate failed: %v", err)
	}
	if err := runCreate("my-token", false); err != nil {
		t.Fatalf("runCreate failed: %v", err)
	}
}

func TestTokenRevoke(t *testing.T) {
	setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{"deleted": true})
	})

	if err := runRevoke(true); err != nil {
		t.Fatalf("runRevoke failed: %v", err)
	}
	if err := runRevoke(false); err != nil {
		t.Fatalf("runRevoke failed: %v", err)
	}
}
