package output

import (
	"encoding/json"
	"fmt"
	"io"
)

// Format represents an output format.
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
	FormatCSV   Format = "csv"
	FormatYAML  Format = "yaml"
)

// ParseFormat parses a format string.
func ParseFormat(s string) Format {
	switch s {
	case "json":
		return FormatJSON
	case "csv":
		return FormatCSV
	case "yaml":
		return FormatYAML
	default:
		return FormatTable
	}
}

// Print outputs data in the specified format.
func Print(w io.Writer, format Format, data interface{}, quiet bool) error {
	switch format {
	case FormatJSON:
		return printJSON(w, data)
	case FormatCSV:
		return printCSV(w, data, quiet)
	case FormatYAML:
		return printYAML(w, data)
	default:
		return printTable(w, data, quiet)
	}
}

// PrintRaw outputs raw JSON from the API in the specified format.
// It automatically unwraps paginated responses (ResultPage with "items" field).
func PrintRaw(w io.Writer, format Format, raw json.RawMessage, quiet bool) error {
	var data interface{}
	if err := json.Unmarshal(raw, &data); err != nil {
		// If we can't parse, just print as-is
		fmt.Fprintln(w, string(raw))
		return nil
	}

	// Unwrap paginated responses: if the object has "items" and "type"=="ResultPage",
	// display the items directly instead of the envelope.
	if m, ok := data.(map[string]interface{}); ok {
		if items, hasItems := m["items"]; hasItems {
			if arr, isArr := items.([]interface{}); isArr {
				if !quiet {
					// Show pagination info
					total, _ := m["total_items"]
					if total != nil {
						fmt.Fprintf(w, "Total: %v items\n\n", total)
					}
				}
				return Print(w, format, arr, quiet)
			}
		}
	}

	return Print(w, format, data, quiet)
}

// PrintItems outputs a slice of raw JSON items.
func PrintItems(w io.Writer, format Format, items []json.RawMessage, quiet bool) error {
	if len(items) == 0 {
		if !quiet {
			fmt.Fprintln(w, "No results found.")
		}
		return nil
	}

	var data []interface{}
	for _, item := range items {
		var v interface{}
		if err := json.Unmarshal(item, &v); err != nil {
			continue
		}
		data = append(data, v)
	}

	return Print(w, format, data, quiet)
}
