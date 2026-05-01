package output

import (
	"encoding/csv"
	"fmt"
	"io"
	"sort"
)

func printCSV(w io.Writer, data interface{}, quiet bool) error {
	items, ok := data.([]interface{})
	if !ok {
		// Single object: print as key,value
		if m, ok := data.(map[string]interface{}); ok {
			items = []interface{}{m}
		} else {
			fmt.Fprintf(w, "%v\n", data)
			return nil
		}
	}

	if len(items) == 0 {
		return nil
	}

	firstMap, ok := items[0].(map[string]interface{})
	if !ok {
		for _, item := range items {
			fmt.Fprintf(w, "%v\n", item)
		}
		return nil
	}

	keys := make([]string, 0, len(firstMap))
	for k := range firstMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	cw := csv.NewWriter(w)

	if !quiet {
		if err := cw.Write(keys); err != nil {
			return err
		}
	}

	for _, item := range items {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		row := make([]string, len(keys))
		for i, k := range keys {
			row[i] = formatValue(m[k])
		}
		if err := cw.Write(row); err != nil {
			return err
		}
	}

	cw.Flush()
	return cw.Error()
}
