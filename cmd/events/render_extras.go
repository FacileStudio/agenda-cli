package events

import (
	"fmt"
	"strconv"

	"github.com/FacileStudio/agenda-cli/internal/caldav"
	"github.com/FacileStudio/agenda-cli/internal/ui"
)

func printExtras(event map[string]any) {
	if tz := text(event["timezone"]); tz != "" {
		ui.PrintKeyValue("Timezone", tz)
	}
	if color := text(event["color"]); color != "" {
		ui.PrintKeyValue("Color", color)
	}
	for _, a := range objects(event["attendees"]) {
		ui.PrintKeyValue("Attendee", formatAttendee(a))
	}
	for _, r := range objects(event["reminders"]) {
		ui.PrintKeyValue("Reminder", fmt.Sprintf("%s min before (%s)", text(r["minutes_before"]), text(r["method"])))
	}
	for _, x := range objects(event["exceptions"]) {
		ui.PrintKeyValue("Exception", formatException(x))
	}
}

func objects(v any) []map[string]any {
	list, _ := v.([]any)
	out := make([]map[string]any, 0, len(list))
	for _, item := range list {
		if m, ok := item.(map[string]any); ok {
			out = append(out, m)
		}
	}
	return out
}

func formatAttendee(a map[string]any) string {
	out := fmt.Sprintf("%s [%s, %s]", text(a["email"]), text(a["role"]), text(a["response"]))
	if name := text(a["name"]); name != "" {
		out = name + " <" + text(a["email"]) + "> [" + text(a["role"]) + ", " + text(a["response"]) + "]"
	}
	return out
}

func formatException(x map[string]any) string {
	date := text(x["date"])
	if x["cancelled"] == true {
		return date + " cancelled"
	}
	out := date + " overridden"
	for _, key := range []string{"title", "location", "start_at", "end_at", "status"} {
		if v := text(x[key]); v != "" {
			out += " " + key + "=" + v
		}
	}
	return out
}

func icsExtras(event map[string]any) *caldav.Extras {
	x := &caldav.Extras{Color: text(event["color"])}
	for _, a := range objects(event["attendees"]) {
		x.Attendees = append(x.Attendees, caldav.Attendee{
			Email: text(a["email"]), Name: text(a["name"]), Role: text(a["role"]), Response: text(a["response"]),
		})
	}
	for _, r := range objects(event["reminders"]) {
		minutes, _ := strconv.Atoi(text(r["minutes_before"]))
		x.Reminders = append(x.Reminders, caldav.Reminder{Minutes: minutes, Method: text(r["method"])})
	}
	for _, e := range objects(event["exceptions"]) {
		if e["cancelled"] == true {
			x.ExDates = append(x.ExDates, parseTime(e["original_start"]))
		}
	}
	return x
}
