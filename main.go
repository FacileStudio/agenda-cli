package main

import (
	"os"

	"github.com/FacileStudio/agenda-cli/cmd"
	"github.com/FacileStudio/agenda-cli/internal/ui"
)

func main() {
	root := cmd.NewRoot()
	if err := root.Execute(); err != nil {
		ui.PrintError(err)
		os.Exit(1)
	}
}
