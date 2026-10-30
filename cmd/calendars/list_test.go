package calendars

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

func TestCalendarsList(t *testing.T) {
	setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/calendars" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": "cal-1", "name": "Work", "color": "#ff0000"},
		})
	})

	if err := runList(true); err != nil {
		t.Fatalf("runList json failed: %v", err)
	}
	if err := runList(false); err != nil {
		t.Fatalf("runList human failed: %v", err)
	}
}

func TestCalendarsShow(t *testing.T) {
	setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/calendars/cal-1" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":       "cal-1",
			"name":     "Work",
			"color":    "#ff0000",
			"owner_id": "usr-1",
		})
	})

	if err := runShow("cal-1", true); err != nil {
		t.Fatalf("runShow json failed: %v", err)
	}
	if err := runShow("cal-1", false); err != nil {
		t.Fatalf("runShow human failed: %v", err)
	}
}

func TestCalendarsDelete(t *testing.T) {
	setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/calendars/cal-2" || r.Method != http.MethodDelete {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{"ok": true})
	})

	if err := runDelete("cal-2", true, true); err != nil {
		t.Fatalf("runDelete json failed: %v", err)
	}
	if err := runDelete("cal-2", true, false); err != nil {
		t.Fatalf("runDelete human failed: %v", err)
	}
}

func TestCalendarsUpdateMergesStoredFields(t *testing.T) {
	var got map[string]any
	setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			json.NewEncoder(w).Encode(map[string]any{
				"id": 1, "name": "Work", "color": "#f00", "description": "keep", "echo_url": "https://e", "role": "owner",
			})
			return
		}
		json.NewDecoder(r.Body).Decode(&got)
		json.NewEncoder(w).Encode(map[string]any{"id": 1})
	})

	cmd := newUpdateCmd()
	cmd.SetArgs([]string{"1", "--color", "#0f0", "--echo-url", "", "--json"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if got["name"] != "Work" || got["description"] != "keep" || got["color"] != "#0f0" || got["echo_url"] != "" {
		t.Errorf("body = %v", got)
	}
	if _, leaked := got["role"]; leaked {
		t.Error("read-only fields must not be sent back")
	}
}

func TestCalendarsCreateSendsEveryField(t *testing.T) {
	var got map[string]any
	setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&got)
		json.NewEncoder(w).Encode(map[string]any{"id": 3, "name": "N"})
	})

	cmd := newCreateCmd()
	cmd.SetArgs([]string{"N", "--color", "#00f", "--description", "d", "--echo-url", "https://e", "--space", "4", "--json"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if got["name"] != "N" || got["description"] != "d" || got["echo_url"] != "https://e" || got["space_id"] != float64(4) {
		t.Errorf("body = %v", got)
	}
	if err := newCreateCmd().Execute(); err == nil {
		t.Error("a calendar with no name must be refused")
	}
}
