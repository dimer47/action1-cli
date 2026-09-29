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

// readDataFlag resolves a --data flag value to its raw bytes. The value is
// inline JSON, @file to read from a file, or - to read from stdin.
func readDataFlag(data string) ([]byte, error) {
	if data == "-" {
		scanner := bufio.NewScanner(os.Stdin)
		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		return []byte(strings.Join(lines, "\n")), nil
	}

	if strings.HasPrefix(data, "@") {
		raw, err := os.ReadFile(data[1:])
		if err != nil {
			return nil, fmt.Errorf("reading file %s: %w", data[1:], err)
		}
		return raw, nil
	}

	return []byte(data), nil
}

// parseDataFlag parses a --data flag value into a JSON object.
func parseDataFlag(data string) (map[string]interface{}, error) {
	raw, err := readDataFlag(data)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}
	return result, nil
}

// parseDataFlagAny parses a --data flag value into any JSON value, keeping the
// payload shape the endpoint expects. Some endpoints take a top-level array
// rather than an object — POST /updates/{orgId}/approvals, for one.
func parseDataFlagAny(data string) (interface{}, error) {
	raw, err := readDataFlag(data)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}
	return result, nil
}
