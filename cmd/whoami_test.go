package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/FacileStudio/agenda-cli/internal/client"
)

func TestWhoamiJSON(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "agenda", "config.yml")
	cfg := &client.Config{
		ServerURL: "https://agenda.example.com",
		Token:     "test-token",
		UserID:    "user-123",
		Email:     "user@example.com",
	}
	if err := cfg.Save(configPath); err != nil {
		t.Fatalf("save config: %v", err)
	}

	t.Setenv("XDG_CONFIG_HOME", tmpDir)
	os.Unsetenv("AGENDA_SERVER_URL")
	os.Unsetenv("AGENDA_TOKEN")

	if err := runWhoami(true); err != nil {
		t.Fatalf("runWhoami json failed: %v", err)
	}
	if err := runWhoami(false); err != nil {
		t.Fatalf("runWhoami human failed: %v", err)
	}
}
