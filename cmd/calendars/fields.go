package calendars

import (
	"github.com/spf13/cobra"
)

type calendarFields struct {
	name        string
	color       string
	description string
	echoURL     string
}

func (f *calendarFields) bind(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.StringVar(&f.name, "name", "", "Calendar name")
	flags.StringVar(&f.color, "color", "", "Calendar color (hex or name)")
	flags.StringVar(&f.description, "description", "", "Calendar description")
	flags.StringVar(&f.echoURL, "echo-url", "", "Echo URL the calendar notifies")
}

func (f *calendarFields) overlay(cmd *cobra.Command, body map[string]any) {
	for _, field := range []struct {
		flag  string
		value string
	}{
		{"name", f.name},
		{"color", f.color},
		{"description", f.description},
		{"echo-url", f.echoURL},
	} {
		if cmd.Flags().Changed(field.flag) {
			body[keyOf(field.flag)] = field.value
		}
	}
}

func keyOf(flag string) string {
	if flag == "echo-url" {
		return "echo_url"
	}
	return flag
}
