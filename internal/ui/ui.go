package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

var (
	successColor = color.New(color.FgGreen).SprintFunc()
	errorColor   = color.New(color.FgRed).SprintFunc()
	warnColor    = color.New(color.FgYellow).SprintFunc()
	infoColor    = color.New(color.FgCyan).SprintFunc()
	dimColor     = color.New(color.Faint).SprintFunc()
)

func PrintSuccess(message string) {
	fmt.Println(successColor("✓"), message)
}

func PrintError(err error) {
	fmt.Fprintln(os.Stderr, errorColor("✗"), err.Error())
}

func PrintWarn(message string) {
	fmt.Println(warnColor("!"), message)
}

func PrintInfo(message string) {
	fmt.Println(infoColor("▸"), message)
}

func PrintDim(message string) {
	fmt.Println(dimColor(message))
}

func PrintKeyValue(key, value string) {
	fmt.Printf("%s %s\n", dimColor(key+":"), value)
}

func Dim(text string) string {
	return dimColor(text)
}

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

func PrintJSON(data interface{}) {
	fmt.Printf("%+v\n", data)
}

func PrintTable(rows [][]string, headers []string) {
	if len(rows) == 0 {
		PrintDim("No results.")
		return
	}

	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	printRow(headers, widths)
	printSeparator(widths)

	for _, row := range rows {
		printRow(row, widths)
	}
}

func printRow(cells []string, widths []int) {
	for i, cell := range cells {
		fmt.Printf("%-*s", widths[i]+2, cell)
	}
	fmt.Println()
}

func printSeparator(widths []int) {
	for _, w := range widths {
		fmt.Print(strings.Repeat("-", w+2))
	}
	fmt.Println()
}
