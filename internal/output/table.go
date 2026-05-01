package output

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"
)

func printTable(w io.Writer, data interface{}, quiet bool) error {
	switch v := data.(type) {
	case []interface{}:
		return printSliceTable(w, v, quiet)
	case map[string]interface{}:
		return printMapTable(w, v, quiet)
	default:
		fmt.Fprintf(w, "%v\n", data)
		return nil
	}
}

func printSliceTable(w io.Writer, items []interface{}, quiet bool) error {
	if len(items) == 0 {
		if !quiet {
			fmt.Fprintln(w, "No results found.")
		}
		return nil
	}

	// Collect all keys from first item
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

	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)

	if !quiet {
		headers := make([]string, len(keys))
		for i, k := range keys {
			headers[i] = strings.ToUpper(k)
		}
		fmt.Fprintln(tw, strings.Join(headers, "\t"))
		sep := make([]string, len(keys))
		for i, h := range headers {
			sep[i] = strings.Repeat("-", len(h))
		}
		fmt.Fprintln(tw, strings.Join(sep, "\t"))
	}

	for _, item := range items {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		vals := make([]string, len(keys))
		for i, k := range keys {
			vals[i] = formatValue(m[k])
		}
		fmt.Fprintln(tw, strings.Join(vals, "\t"))
	}

	return tw.Flush()
}

func printMapTable(w io.Writer, m map[string]interface{}, quiet bool) error {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if !quiet {
			fmt.Fprintf(tw, "%s\t%s\n", k, formatValue(m[k]))
		} else {
			fmt.Fprintf(tw, "%s\n", formatValue(m[k]))
		}
	}

	return tw.Flush()
}

func formatValue(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case float64:
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}
		return fmt.Sprintf("%.2f", val)
	case bool:
		if val {
			return "true"
		}
		return "false"
	case map[string]interface{}, []interface{}:
		return "[...]"
	default:
		return fmt.Sprintf("%v", v)
	}
}
