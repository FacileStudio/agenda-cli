package events

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

func TestEventsList(t *testing.T) {
	setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/calendars/cal-1/events" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": "ev-1", "title": "Meeting", "start": "2026-10-27T10:00:00Z", "end": "2026-10-27T11:00:00Z"},
		})
	})

	filterJSON := listFilter{calendarID: "cal-1", jsonOutput: true}
	if err := runList(filterJSON); err != nil {
		t.Fatalf("runList json failed: %v", err)
	}

	filterHuman := listFilter{calendarID: "cal-1", jsonOutput: false}
	if err := runList(filterHuman); err != nil {
		t.Fatalf("runList human failed: %v", err)
	}
}

func TestEventsShow(t *testing.T) {
	setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/events/ev-1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":          "ev-1",
			"title":       "Meeting",
			"start":       "2026-10-27T10:00:00Z",
			"end":         "2026-10-27T11:00:00Z",
			"calendar_id": "cal-1",
		})
	})

	if err := runShow("ev-1", false, true); err != nil {
		t.Fatalf("runShow json failed: %v", err)
	}
	if err := runShow("ev-1", false, false); err != nil {
		t.Fatalf("runShow human failed: %v", err)
	}
}
