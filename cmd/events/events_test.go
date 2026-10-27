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

func TestEventsCreate(t *testing.T) {
	setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/calendars/cal-1/events" || r.Method != http.MethodPost {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if payload["start_at"] != "2026-10-27T10:00:00Z" {
			t.Errorf("expected start_at in payload, got %v", payload["start_at"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":       "ev-2",
			"title":    "New Event",
			"start_at": "2026-10-27T10:00:00Z",
			"end_at":   "2026-10-27T11:00:00Z",
		})
	})

	opts := createEventOptions{
		calendarID: "cal-1",
		title:      "New Event",
		start:      "2026-10-27T10:00:00Z",
		end:        "2026-10-27T11:00:00Z",
	}
	if err := runCreate(opts, true); err != nil {
		t.Fatalf("runCreate json failed: %v", err)
	}
	if err := runCreate(opts, false); err != nil {
		t.Fatalf("runCreate human failed: %v", err)
	}
}

func TestEventsUpdate(t *testing.T) {
	setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			json.NewEncoder(w).Encode(map[string]any{
				"id":       "ev-2",
				"title":    "Old Event",
				"start_at": "2026-10-27T10:00:00Z",
				"end_at":   "2026-10-27T11:00:00Z",
			})
			return
		}
		if r.URL.Path != "/events/ev-2" || r.Method != http.MethodPut {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":       "ev-2",
			"title":    "Updated Event",
			"start_at": "2026-10-27T10:00:00Z",
			"end_at":   "2026-10-27T11:00:00Z",
		})
	})

	opts := updateEventOptions{title: "Updated Event"}
	if err := runUpdate("ev-2", opts, true); err != nil {
		t.Fatalf("runUpdate json failed: %v", err)
	}
	if err := runUpdate("ev-2", opts, false); err != nil {
		t.Fatalf("runUpdate human failed: %v", err)
	}
}

func TestEventsDelete(t *testing.T) {
	setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/events/ev-2" || r.Method != http.MethodDelete {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{"ok": true})
	})

	if err := runDelete("ev-2", true, true); err != nil {
		t.Fatalf("runDelete json failed: %v", err)
	}
	if err := runDelete("ev-2", true, false); err != nil {
		t.Fatalf("runDelete human failed: %v", err)
	}
}
