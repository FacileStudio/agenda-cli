package caldav

import (
	"fmt"
	"net/http"

	"github.com/FacileStudio/agenda-cli/internal/ui"
	"github.com/spf13/cobra"
)

func newTestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "test",
		Short: "Test CalDAV connectivity",
		Long:  "Perform a PROPFIND probe against the CalDAV home set to verify auth and server connectivity.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTest()
		},
	}
}

func runTest() error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	caldavURL := fmt.Sprintf("%s/dav/%s/calendars", config.ServerURL, config.Email)

	ui.PrintInfo("Probing CalDAV endpoint...")
	ui.PrintDim(caldavURL)

	req, err := http.NewRequest("PROPFIND", caldavURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", config.AuthHeader())
	req.Header.Set("Depth", "0")
	req.Header.Set("Content-Type", "application/xml")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("CalDAV test failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 207 || resp.StatusCode == 200 {
		ui.PrintSuccess("CalDAV connection successful")
		return nil
	}

	return fmt.Errorf("CalDAV test failed with status %d", resp.StatusCode)
}
