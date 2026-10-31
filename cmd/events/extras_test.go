package events

import (
	"reflect"
	"strings"
	"testing"
)

func TestCreateSendsAttendeesRemindersAndColor(t *testing.T) {
	var got map[string]any
	captureBody(t, &got, map[string]any{"id": 2, "title": "T"}, nil)
	err := execCmd(t, newCreateCmd(), "--calendar", "1", "--title", "T", "--start", "2026-11-02", "--end", "2026-11-03",
		"--attendee", "bob@example.com,name=Bob,role=optional,response=accepted", "--attendee", "eve@example.com",
		"--reminder", "15", "--reminder", "5:audio", "--color", "#ff8800", "--json")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	wantAttendees := []any{
		map[string]any{"email": "bob@example.com", "name": "Bob", "role": "optional", "response": "accepted"},
		map[string]any{"email": "eve@example.com"},
	}
	wantReminders := []any{
		map[string]any{"minutes_before": float64(15), "method": "display"},
		map[string]any{"minutes_before": float64(5), "method": "audio"},
	}
	if !reflect.DeepEqual(got["attendees"], wantAttendees) || !reflect.DeepEqual(got["reminders"], wantReminders) {
		t.Errorf("lists = %v / %v", got["attendees"], got["reminders"])
	}
	if got["color"] != "#ff8800" {
		t.Errorf("color = %v", got["color"])
	}
}

func TestCreateRejectsBadExtrasBeforeTheRequest(t *testing.T) {
	base := []string{"--calendar", "1", "--title", "T", "--start", "2026-11-02", "--end", "2026-11-03"}
	for name, extra := range map[string][]string{
		"email":    {"--attendee", "nobody"},
		"role":     {"--attendee", "a@example.com,role=boss"},
		"option":   {"--attendee", "a@example.com,mood=happy"},
		"reminder": {"--reminder", "soon"},
		"method":   {"--reminder", "5:pigeon"},
		"color":    {"--color", "red"},
		"conflict": {"--attendee", "a@example.com", "--clear-attendees"},
	} {
		if err := execCmd(t, newCreateCmd(), append(append([]string{}, base...), extra...)...); err == nil {
			t.Errorf("%s: a bad value must be refused", name)
		}
	}
}

func TestUpdateKeepsZoneAndOmitsUntouchedLists(t *testing.T) {
	var got map[string]any
	existing := map[string]any{
		"id": 2, "title": "Old", "start_at": "2026-11-02T08:00:00Z", "end_at": "2026-11-02T08:30:00Z",
		"recurrence_rule": "FREQ=DAILY", "timezone": "Europe/Paris", "attendees": []any{map[string]any{"email": "a@b.co"}},
	}
	captureBody(t, &got, map[string]any{"id": 2}, existing)
	err := execCmd(t, newUpdateCmd(), "2", "--start", "2026-11-03 09:00", "--end", "2026-11-03 09:30", "--json")
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if got["timezone"] != "Europe/Paris" || got["start_at"] != "2026-11-03T08:00:00Z" {
		t.Errorf("the stored zone was not kept: %v", got)
	}
	for _, key := range []string{"attendees", "reminders", "color"} {
		if _, sent := got[key]; sent {
			t.Errorf("%s was sent although no flag touched it", key)
		}
	}
}

func TestUpdateClearsLists(t *testing.T) {
	var got map[string]any
	captureBody(t, &got, map[string]any{"id": 2}, map[string]any{"id": 2, "title": "Old"})
	err := execCmd(t, newUpdateCmd(), "2", "--clear-attendees", "--clear-reminders", "--color", "", "--json")
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if a, ok := got["attendees"].([]any); !ok || len(a) != 0 {
		t.Errorf("attendees = %v, want []", got["attendees"])
	}
	if r, ok := got["reminders"].([]any); !ok || len(r) != 0 || got["color"] != "" {
		t.Errorf("reminders/color = %v / %v", got["reminders"], got["color"])
	}
}

func TestShowRendersExtrasAndICS(t *testing.T) {
	event := map[string]any{
		"uid": "u@agenda", "title": "T", "start_at": "2026-11-02T08:00:00Z", "end_at": "2026-11-02T09:00:00Z",
		"color":      "#112233",
		"attendees":  []any{map[string]any{"email": "a@b.co", "name": "Ann", "role": "optional", "response": "accepted"}},
		"reminders":  []any{map[string]any{"minutes_before": float64(10), "method": "audio"}},
		"exceptions": []any{map[string]any{"date": "2026-11-04", "cancelled": true, "original_start": "2026-11-04T08:00:00Z"}},
	}
	data, err := renderICS(event)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"ATTENDEE;", "PARTSTAT=ACCEPTED", "BEGIN:VALARM", "TRIGGER:-PT10M", "EXDATE:20261104T080000Z", "COLOR:#112233"} {
		if !strings.Contains(data, want) {
			t.Errorf("ICS lacks %q:\n%s", want, data)
		}
	}
	if got := formatException(objects(event["exceptions"])[0]); got != "2026-11-04 cancelled" {
		t.Errorf("formatException = %q", got)
	}
	printEvent(event)
}
