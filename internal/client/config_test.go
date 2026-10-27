package client

import (
	"os"
	"testing"
)

func TestConfigLoadWithEnvVars(t *testing.T) {
	os.Setenv("AGENDA_SERVER_URL", "https://agenda.test.example")
	os.Setenv("AGENDA_TOKEN", "test-token-123")
	defer func() {
		os.Unsetenv("AGENDA_SERVER_URL")
		os.Unsetenv("AGENDA_TOKEN")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.ServerURL != "https://agenda.test.example" {
		t.Errorf("expected ServerURL https://agenda.test.example, got %s", cfg.ServerURL)
	}
	if cfg.Token != "test-token-123" {
		t.Errorf("expected Token test-token-123, got %s", cfg.Token)
	}
	if err := cfg.RequireAuth(); err != nil {
		t.Errorf("expected RequireAuth to succeed, got %v", err)
	}
}

func TestConfigRequireAuthMissing(t *testing.T) {
	cfg := &Config{}
	if err := cfg.RequireAuth(); err == nil {
		t.Error("expected error for empty config, got nil")
	}

	cfg.ServerURL = "https://agenda.test.example"
	if err := cfg.RequireAuth(); err == nil {
		t.Error("expected error for missing token, got nil")
	}
}
