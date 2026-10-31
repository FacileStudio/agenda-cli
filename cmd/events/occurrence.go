package events

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/FacileStudio/agenda-cli/internal/client"
	"github.com/spf13/cobra"
)

func occurrencePath(id, when, suffix string) string {
	return "/events/" + id + "/occurrences/" + url.PathEscape(when) + suffix
}

func newCancelOccurrenceCmd() *cobra.Command {
	return newOccurrenceCmd(occurrenceAction{
		use:     "cancel-occurrence <id>",
		short:   "Cancel one occurrence of a recurring event",
		long:    "Exclude a single occurrence from a recurring event (an EXDATE). The rest of the series is untouched.",
		message: "Occurrence cancelled",
		send: func(api *client.Client, id, when string) (*http.Response, error) {
			return api.Delete(occurrencePath(id, when, ""))
		},
	})
}

func newRestoreOccurrenceCmd() *cobra.Command {
	return newOccurrenceCmd(occurrenceAction{
		use:     "restore-occurrence <id>",
		short:   "Restore a cancelled or overridden occurrence",
		long:    "Undo cancel-occurrence or update --occurrence so the occurrence follows the series again.",
		message: "Occurrence restored",
		send: func(api *client.Client, id, when string) (*http.Response, error) {
			return api.Post(occurrencePath(id, when, "/restore"), nil)
		},
	})
}

type occurrenceAction struct {
	use     string
	short   string
	long    string
	message string
	send    func(api *client.Client, id, when string) (*http.Response, error)
}

func newOccurrenceCmd(action occurrenceAction) *cobra.Command {
	var (
		date       string
		jsonOutput bool
	)
	cmd := &cobra.Command{
		Use:   action.use,
		Short: action.short,
		Long:  action.long + "\n\n--date is the occurrence day (YYYY-MM-DD, in the event's timezone) or its start as RFC3339.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runOccurrence(action, args[0], date, jsonOutput)
		},
	}
	cmd.Flags().StringVar(&date, "date", "", "Occurrence date YYYY-MM-DD, or RFC3339 start")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	cmd.MarkFlagRequired("date")
	return cmd
}

func runOccurrence(action occurrenceAction, id, date string, jsonOutput bool) error {
	if date == "" {
		return fmt.Errorf("--date is required")
	}
	api, err := loadAPI()
	if err != nil {
		return err
	}
	resp, err := action.send(api, id, date)
	if err != nil {
		return err
	}
	return printResult(api, resp, action.message, jsonOutput)
}
