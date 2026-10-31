package events

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var validStatuses = []string{"confirmed", "tentative", "cancelled"}

type eventFields struct {
	title              string
	description        string
	location           string
	start              string
	end                string
	recurrence         string
	status             string
	conferenceURL      string
	conferenceProvider string
	timezone           string
	allDay             bool
	extras             eventExtras
}

type textField struct {
	flag  string
	key   string
	value *string
}

func (f *eventFields) bind(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.StringVar(&f.title, "title", "", "Event title")
	flags.StringVar(&f.start, "start", "", "Start (RFC3339, or YYYY-MM-DD[ HH:MM] read in --timezone, UTC by default)")
	flags.StringVar(&f.end, "end", "", "End, same formats as --start")
	flags.BoolVar(&f.allDay, "all-day", false, "All-day event, --all-day=false to clear it")
	flags.StringVar(&f.location, "location", "", "Event location")
	flags.StringVar(&f.description, "description", "", "Event description")
	flags.StringVar(&f.recurrence, "recurrence", "", `RFC 5545 RRULE such as "FREQ=WEEKLY;BYDAY=MO", empty clears it`)
	flags.StringVar(&f.status, "status", "", "confirmed, tentative or cancelled")
	flags.StringVar(&f.conferenceURL, "conference-url", "", "Video call or meeting link")
	flags.StringVar(&f.conferenceProvider, "conference-provider", "", "Label for the link, such as Zoom")
	flags.StringVar(&f.timezone, "timezone", "", "IANA zone such as Europe/Paris: reads naive --start/--end and anchors a recurrence")
	f.extras.bind(cmd)
}

func (f *eventFields) textFields() []textField {
	return []textField{
		{"title", "title", &f.title},
		{"description", "description", &f.description},
		{"location", "location", &f.location},
		{"recurrence", "recurrence_rule", &f.recurrence},
		{"status", "status", &f.status},
		{"conference-url", "conference_url", &f.conferenceURL},
		{"conference-provider", "conference_provider", &f.conferenceProvider},
		{"timezone", "timezone", &f.timezone},
	}
}

func (f *eventFields) overlay(cmd *cobra.Command, body map[string]any) error {
	changed := cmd.Flags().Changed
	if changed("status") && !slices.Contains(validStatuses, f.status) {
		return fmt.Errorf("invalid --status %q: use %s", f.status, strings.Join(validStatuses, ", "))
	}
	loc, err := loadZone(f.timezone)
	if err != nil {
		return err
	}
	for _, field := range f.textFields() {
		if changed(field.flag) {
			body[field.key] = *field.value
		}
	}
	if changed("all-day") {
		body["is_all_day"] = f.allDay
	}
	if err := f.extras.overlay(cmd, body); err != nil {
		return err
	}
	return f.overlayTimes(cmd, body, loc)
}

func (f *eventFields) overlayTimes(cmd *cobra.Command, body map[string]any, loc *time.Location) error {
	for _, t := range []struct{ flag, key, raw string }{
		{"start", "start_at", f.start},
		{"end", "end_at", f.end},
	} {
		if !cmd.Flags().Changed(t.flag) {
			continue
		}
		parsed, err := parseDateTime(t.raw, loc)
		if err != nil {
			return fmt.Errorf("invalid --%s: %w", t.flag, err)
		}
		body[t.key] = parsed
	}
	return nil
}

func (f *eventFields) adoptStoredZone(cmd *cobra.Command, existing map[string]any) {
	if cmd.Flags().Changed("timezone") {
		return
	}
	if stored, ok := existing["timezone"].(string); ok {
		f.timezone = stored
	}
}

func loadZone(name string) (*time.Location, error) {
	if name == "" {
		return time.UTC, nil
	}
	loc, err := time.LoadLocation(name)
	if err != nil {
		return nil, fmt.Errorf("invalid --timezone %q: use an IANA name such as Europe/Paris", name)
	}
	return loc, nil
}

func parseDateTime(raw string, loc *time.Location) (string, error) {
	raw = strings.TrimSpace(raw)
	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		return t.UTC().Format(time.RFC3339), nil
	}
	for _, layout := range []string{"2006-01-02T15:04:05", "2006-01-02 15:04:05", "2006-01-02T15:04", "2006-01-02 15:04", "2006-01-02"} {
		if t, err := time.ParseInLocation(layout, raw, loc); err == nil {
			return t.UTC().Format(time.RFC3339), nil
		}
	}
	return "", fmt.Errorf("expected RFC3339 or YYYY-MM-DD HH:MM")
}
