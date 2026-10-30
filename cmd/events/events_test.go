package events

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/spf13/cobra"
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

func execCmd(t *testing.T, cmd *cobra.Command, args ...string) error {
	t.Helper()
	cmd.SetArgs(args)
	return cmd.Execute()
}

func captureBody(t *testing.T, got *map[string]any, reply map[string]any, existing map[string]any) {
	t.Helper()
	setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			json.NewEncoder(w).Encode(existing)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(got); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		json.NewEncoder(w).Encode(reply)
	})
}

func TestEventsCreateSendsEveryField(t *testing.T) {
	var got map[string]any
	captureBody(t, &got, map[string]any{"id": 2, "title": "T"}, nil)

	err := execCmd(t, newCreateCmd(), "--calendar", "1", "--title", "T", "--start", "2026-10-27 09:00",
		"--end", "2026-10-27 10:00", "--timezone", "Europe/Paris", "--recurrence", "FREQ=WEEKLY;BYDAY=MO",
		"--status", "tentative", "--conference-url", "https://meet.example/x", "--conference-provider", "Jitsi",
		"--location", "Room", "--description", "D", "--json")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	want := map[string]any{
		"title": "T", "start_at": "2026-10-27T08:00:00Z", "end_at": "2026-10-27T09:00:00Z",
		"timezone": "Europe/Paris", "recurrence_rule": "FREQ=WEEKLY;BYDAY=MO", "status": "tentative",
		"conference_url": "https://meet.example/x", "conference_provider": "Jitsi", "location": "Room", "description": "D",
	}
	for key, value := range want {
		if got[key] != value {
			t.Errorf("%s = %v, want %v", key, got[key], value)
		}
	}
}

func TestEventsCreateRejectsBadStatusAndZone(t *testing.T) {
	base := []string{"--calendar", "1", "--title", "T", "--start", "2026-10-27", "--end", "2026-10-28"}
	if err := execCmd(t, newCreateCmd(), append(base, "--status", "busy")...); err == nil {
		t.Error("an unknown status must be refused before the request")
	}
	if err := execCmd(t, newCreateCmd(), append(base, "--timezone", "Mars/Base")...); err == nil {
		t.Error("an unknown timezone must be refused before the request")
	}
}

func TestEventsUpdateKeepsStoredFieldsAndClears(t *testing.T) {
	var got map[string]any
	existing := map[string]any{
		"id": 2, "title": "Old", "start_at": "2026-10-27T10:00:00Z", "end_at": "2026-10-27T11:00:00Z",
		"is_all_day": true, "location": "Here", "recurrence_rule": "FREQ=DAILY", "status": "tentative",
		"conference_url": "https://x", "created_at": "ignored",
	}
	captureBody(t, &got, map[string]any{"id": 2, "title": "New"}, existing)

	err := execCmd(t, newUpdateCmd(), "2", "--title", "New", "--location", "", "--all-day=false", "--calendar", "7", "--json")
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if got["recurrence_rule"] != "FREQ=DAILY" || got["status"] != "tentative" || got["conference_url"] != "https://x" {
		t.Errorf("stored fields were not carried: %v", got)
	}
	if got["location"] != "" || got["is_all_day"] != false || got["title"] != "New" || got["calendar_id"] != float64(7) {
		t.Errorf("overrides were not applied: %v", got)
	}
	if _, leaked := got["created_at"]; leaked {
		t.Error("read-only fields must not be sent back")
	}
}

func TestEventsUpdateNeedsAField(t *testing.T) {
	if err := execCmd(t, newUpdateCmd(), "2", "--json"); err == nil {
		t.Error("an update with no field flag must fail")
	}
}

func TestEventsListSearchIsLocal(t *testing.T) {
	setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("tz") != "Europe/Paris" {
			t.Errorf("tz = %q", r.URL.Query().Get("tz"))
		}
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 1, "title": "Standup"}, {"id": 2, "title": "Lunch", "location": "STAND bar"}, {"id": 3, "title": "Other"},
		})
	})
	got := filterEvents([]map[string]any{{"title": "Standup"}, {"title": "x", "location": "STAND"}, {"title": "no"}}, "stand")
	if len(got) != 2 {
		t.Fatalf("filterEvents kept %d, want 2", len(got))
	}
	if err := runList(listFilter{calendarID: "1", search: "stand", timezone: "Europe/Paris", jsonOutput: true}); err != nil {
		t.Fatal(err)
	}
}

func TestEventsDeleteReportsServerErrors(t *testing.T) {
	setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	})
	if err := runDelete("9", true, true); err == nil {
		t.Error("a 404 on delete must be an error")
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
