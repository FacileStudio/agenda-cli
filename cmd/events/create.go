package events

import (
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	var (
		fields     eventFields
		calendarID string
		jsonOutput bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new event",
		Long:  "Create a new event in a calendar.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreate(cmd, &fields, calendarID, jsonOutput)
		},
	}

	fields.bind(cmd)
	cmd.Flags().StringVar(&calendarID, "calendar", "", "Calendar ID")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")

	cmd.MarkFlagRequired("calendar")
	cmd.MarkFlagRequired("title")
	cmd.MarkFlagRequired("start")
	cmd.MarkFlagRequired("end")

	return cmd
}

func runCreate(cmd *cobra.Command, fields *eventFields, calendarID string, jsonOutput bool) error {
	body := map[string]any{}
	if err := fields.overlay(cmd, body); err != nil {
		return err
	}

	api, err := loadAPI()
	if err != nil {
		return err
	}

	resp, err := api.Post("/calendars/"+calendarID+"/events", body)
	if err != nil {
		return err
	}
	return printResult(api, resp, "Event created", jsonOutput)
}
