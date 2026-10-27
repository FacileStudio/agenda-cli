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

func TestCalendarsCreate(t *testing.T) {
	setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/calendars" || r.Method != http.MethodPost {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":    "cal-2",
			"name":  "Project",
			"color": "#00ff00",
		})
	})

	if err := runCreate("Project", "#00ff00", true); err != nil {
		t.Fatalf("runCreate json failed: %v", err)
	}
	if err := runCreate("Project", "#00ff00", false); err != nil {
		t.Fatalf("runCreate human failed: %v", err)
	}
}

func TestCalendarsUpdate(t *testing.T) {
	setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/calendars/cal-2" || r.Method != http.MethodPut {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":    "cal-2",
			"name":  "Renamed",
			"color": "#0000ff",
		})
	})

	if err := runUpdate("cal-2", "Renamed", "#0000ff", true); err != nil {
		t.Fatalf("runUpdate json failed: %v", err)
	}
	if err := runUpdate("cal-2", "Renamed", "#0000ff", false); err != nil {
		t.Fatalf("runUpdate human failed: %v", err)
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
