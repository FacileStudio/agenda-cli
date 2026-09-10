package ui

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var (
	successColor = color.New(color.FgGreen).SprintFunc()
	errorColor   = color.New(color.FgRed).SprintFunc()
	warnColor    = color.New(color.FgYellow).SprintFunc()
	infoColor    = color.New(color.FgCyan).SprintFunc()
	dimColor     = color.New(color.Faint).SprintFunc()
)

// PrintSuccess prints a green-check marked success message.
func PrintSuccess(message string) {
	fmt.Println(successColor("✓"), message)
}

// PrintError prints a red-cross marked error message on stderr.
func PrintError(err error) {
	fmt.Fprintln(os.Stderr, errorColor("✗"), err.Error())
}

// PrintWarn prints a yellow-bang marked warning message.
func PrintWarn(message string) {
	fmt.Println(warnColor("!"), message)
}

// PrintInfo prints a cyan-arrow info message.
func PrintInfo(message string) {
	fmt.Println(infoColor("▸"), message)
}

// PrintDim prints a faint message.
func PrintDim(message string) {
	fmt.Println(dimColor(message))
}

// PrintKeyValue prints a dim key followed by its value.
func PrintKeyValue(key, value string) {
	fmt.Printf("%s %s\n", dimColor(key+":"), value)
}

// Init configures colouring based on the NO_COLOR env var and terminal detection.
func Init() {
	color.NoColor = true
	if os.Getenv("NO_COLOR") == "" && isTerminal(os.Stderr) {
		color.NoColor = false
	}
}

func isTerminal(f *os.File) bool {
	fileInfo, err := f.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
