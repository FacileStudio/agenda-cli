package events

import (
	"encoding/json"
	"net/http"
	"testing"
)

func recordRequest(t *testing.T, method, path *string, body *map[string]any) {
	t.Helper()
	setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			json.NewEncoder(w).Encode(map[string]any{"id": 2, "timezone": "Europe/Paris"})
			return
		}
		*method, *path = r.Method, r.URL.EscapedPath()
		*body = nil
		json.NewDecoder(r.Body).Decode(body)
		json.NewEncoder(w).Encode(map[string]any{"id": 2, "title": "T"})
	})
}

func TestCancelAndRestoreOccurrence(t *testing.T) {
	var method, path string
	var body map[string]any
	recordRequest(t, &method, &path, &body)

	if err := execCmd(t, newCancelOccurrenceCmd(), "2", "--date", "2026-11-04", "--json"); err != nil {
		t.Fatal(err)
	}
	if method != http.MethodDelete || path != "/events/2/occurrences/2026-11-04" {
		t.Errorf("cancel = %s %s", method, path)
	}
	if err := execCmd(t, newRestoreOccurrenceCmd(), "2", "--date", "2026-11-04T08:00:00Z", "--json"); err != nil {
		t.Fatal(err)
	}
	if method != http.MethodPost || path != "/events/2/occurrences/2026-11-04T08:00:00Z/restore" {
		t.Errorf("restore = %s %s", method, path)
	}
}

func TestOccurrenceNeedsADate(t *testing.T) {
	if err := execCmd(t, newCancelOccurrenceCmd(), "2"); err == nil {
		t.Error("cancel-occurrence without --date must fail")
	}
}

func TestUpdateOccurrenceSendsOnlyOccurrenceFields(t *testing.T) {
	var method, path string
	var body map[string]any
	recordRequest(t, &method, &path, &body)

	err := execCmd(t, newUpdateCmd(), "2", "--occurrence", "2026-11-03", "--title", "Moved",
		"--start", "2026-11-03 14:00", "--location", "Room 2", "--json")
	if err != nil {
		t.Fatal(err)
	}
	if method != http.MethodPut || path != "/events/2/occurrences/2026-11-03" {
		t.Fatalf("request = %s %s", method, path)
	}
	if body["title"] != "Moved" || body["location"] != "Room 2" || body["start_at"] != "2026-11-03T13:00:00Z" {
		t.Errorf("body = %v", body)
	}
	if _, leaked := body["timezone"]; leaked || len(body) != 3 {
		t.Errorf("series fields leaked into the occurrence body: %v", body)
	}
}

func TestUpdateOccurrenceRefusesSeriesFlags(t *testing.T) {
	var method, path string
	var body map[string]any
	recordRequest(t, &method, &path, &body)
	for _, flag := range [][]string{{"--recurrence", "FREQ=DAILY"}, {"--attendee", "a@b.co"}, {"--color", "#000000"}, {"--status", "cancelled"}} {
		args := append([]string{"2", "--occurrence", "2026-11-03", "--title", "x"}, flag...)
		if err := execCmd(t, newUpdateCmd(), args...); err == nil {
			t.Errorf("%v must be refused with --occurrence", flag)
		}
	}
	if err := execCmd(t, newUpdateCmd(), "2", "--occurrence", "2026-11-03"); err == nil {
		t.Error("an occurrence update with no field must fail")
	}
}
