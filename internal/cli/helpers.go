package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// confirmAction prompts the user for confirmation.
func confirmAction(action string) bool {
	fmt.Printf("Are you sure you want to %s? [y/N]: ", action)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
		return answer == "y" || answer == "yes"
	}
	return false
}

// rawToInterface converts a slice of json.RawMessage to []interface{}.
func rawToInterface(items []json.RawMessage) []interface{} {
	result := make([]interface{}, 0, len(items))
	for _, item := range items {
		var v interface{}
		if err := json.Unmarshal(item, &v); err == nil {
			result = append(result, v)
		}
	}
	return result
}

// parseDataFlag parses a --data flag value which can be inline JSON, @file, or -.
func parseDataFlag(data string) (map[string]interface{}, error) {
	var raw []byte

	if data == "-" {
		scanner := bufio.NewScanner(os.Stdin)
		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		raw = []byte(strings.Join(lines, "\n"))
	} else if strings.HasPrefix(data, "@") {
		var err error
		raw, err = os.ReadFile(data[1:])
		if err != nil {
			return nil, fmt.Errorf("reading file %s: %w", data[1:], err)
		}
	} else {
		raw = []byte(data)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}
	return result, nil
}
