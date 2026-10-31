package events

import (
	"fmt"
	"net/mail"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	colorPattern    = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)
	validRoles      = []string{"required", "optional"}
	validResponses  = []string{"needs-action", "accepted", "declined", "tentative"}
	validMethods    = []string{"display", "audio", "email"}
	attendeeKeyHelp = "name=, role=required|optional, response=needs-action|accepted|declined|tentative"
)

type eventExtras struct {
	attendees      []string
	reminders      []string
	color          string
	clearAttendees bool
	clearReminders bool
}

func (x *eventExtras) bind(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.StringArrayVar(&x.attendees, "attendee", nil,
		"Attendee EMAIL[,"+attendeeKeyHelp+"], repeatable; replaces the stored list")
	flags.BoolVar(&x.clearAttendees, "clear-attendees", false, "Remove every attendee")
	flags.StringArrayVar(&x.reminders, "reminder", nil,
		"Reminder MINUTES[:display|audio|email] before the start, repeatable; replaces the stored list")
	flags.BoolVar(&x.clearReminders, "clear-reminders", false, "Remove every reminder")
	flags.StringVar(&x.color, "color", "", "Event color #RRGGBB, empty clears it")
}

func (x *eventExtras) overlay(cmd *cobra.Command, body map[string]any) error {
	changed := cmd.Flags().Changed
	if changed("attendee") && x.clearAttendees || changed("reminder") && x.clearReminders {
		return fmt.Errorf("--clear-attendees and --clear-reminders cannot be combined with --attendee and --reminder")
	}
	if changed("color") {
		if x.color != "" && !colorPattern.MatchString(x.color) {
			return fmt.Errorf("invalid --color %q: use #RRGGBB", x.color)
		}
		body["color"] = x.color
	}
	return x.overlayLists(body)
}

func (x *eventExtras) overlayLists(body map[string]any) error {
	if x.clearAttendees {
		body["attendees"] = []map[string]string{}
	}
	if x.clearReminders {
		body["reminders"] = []map[string]any{}
	}
	if len(x.attendees) > 0 {
		list, err := parseAttendees(x.attendees)
		if err != nil {
			return err
		}
		body["attendees"] = list
	}
	if len(x.reminders) > 0 {
		list, err := parseReminders(x.reminders)
		if err != nil {
			return err
		}
		body["reminders"] = list
	}
	return nil
}

func parseAttendees(specs []string) ([]map[string]string, error) {
	out := make([]map[string]string, 0, len(specs))
	for _, spec := range specs {
		a, err := parseAttendee(spec)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, nil
}

func parseAttendee(spec string) (map[string]string, error) {
	parts := strings.Split(spec, ",")
	addr, err := mail.ParseAddress(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, fmt.Errorf("invalid --attendee %q: start with an email address", spec)
	}
	out := map[string]string{"email": addr.Address}
	for _, part := range parts[1:] {
		key, value, ok := strings.Cut(part, "=")
		if !ok || !attendeeValueOK(strings.TrimSpace(key), strings.TrimSpace(value)) {
			return nil, fmt.Errorf("invalid --attendee option %q: use %s", part, attendeeKeyHelp)
		}
		out[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}
	return out, nil
}

func attendeeValueOK(key, value string) bool {
	switch key {
	case "name":
		return true
	case "role":
		return slices.Contains(validRoles, value)
	case "response":
		return slices.Contains(validResponses, value)
	}
	return false
}

func parseReminders(specs []string) ([]map[string]any, error) {
	out := make([]map[string]any, 0, len(specs))
	for _, spec := range specs {
		minutes, method, _ := strings.Cut(spec, ":")
		n, err := strconv.Atoi(strings.TrimSpace(minutes))
		if err != nil || n < 0 {
			return nil, fmt.Errorf("invalid --reminder %q: use MINUTES[:display|audio|email]", spec)
		}
		if method == "" {
			method = "display"
		}
		if !slices.Contains(validMethods, method) {
			return nil, fmt.Errorf("invalid --reminder method %q: use %s", method, strings.Join(validMethods, ", "))
		}
		out = append(out, map[string]any{"minutes_before": n, "method": method})
	}
	return out, nil
}
