package events

import (
	"fmt"
	"strings"

	"github.com/FacileStudio/agenda-cli/internal/client"
	"github.com/FacileStudio/agenda-cli/internal/ui"
	"net/http"
)

func printResult(api *client.Client, resp *http.Response, message string, jsonOutput bool) error {
	var event map[string]any
	if err := api.DecodeJSON(resp, &event); err != nil {
		return err
	}
	if jsonOutput {
		return ui.PrintJSON(event)
	}
	ui.PrintSuccess(message)
	ui.PrintKeyValue("ID", text(event["id"]))
	ui.PrintKeyValue("Title", text(event["title"]))
	return nil
}

func text(v any) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

func printEvent(event map[string]any) {
	rows := [][2]string{
		{"ID", text(event["id"])},
		{"Title", text(event["title"])},
		{"Start", text(event["start_at"])},
		{"End", text(event["end_at"])},
		{"All Day", text(event["is_all_day"])},
		{"Recurrence", text(event["recurrence_rule"])},
		{"Status", text(event["status"])},
		{"Location", text(event["location"])},
		{"Conference URL", text(event["conference_url"])},
		{"Conference Provider", text(event["conference_provider"])},
		{"Description", text(event["description"])},
		{"Calendar", text(event["calendar_id"])},
	}
	for _, row := range rows {
		ui.PrintKeyValue(row[0], row[1])
	}
}

func printEvents(events []map[string]any) {
	ui.PrintSuccess(fmt.Sprintf("Found %d events", len(events)))
	for _, event := range events {
		ui.PrintKeyValue("ID", text(event["id"]))
		ui.PrintKeyValue("Title", text(event["title"]))
		ui.PrintKeyValue("Start", text(event["start_at"]))
		ui.PrintKeyValue("End", text(event["end_at"]))
		if rule := text(event["recurrence_rule"]); rule != "" {
			ui.PrintKeyValue("Recurrence", rule)
		}
		fmt.Println()
	}
}

func matchesSearch(event map[string]any, needle string) bool {
	needle = strings.ToLower(needle)
	for _, key := range []string{"title", "description", "location"} {
		if strings.Contains(strings.ToLower(text(event[key])), needle) {
			return true
		}
	}
	return false
}
