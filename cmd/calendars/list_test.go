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
