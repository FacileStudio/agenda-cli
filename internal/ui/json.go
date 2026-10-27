package ui

import (
	"encoding/json"
	"fmt"
)

// PrintJSON formats data as indented JSON and prints it.
func PrintJSON(data any) error {
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}
