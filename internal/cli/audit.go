package cli

import (
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

func newAuditCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "audit",
		Short: "Manage audit trail",
	}

	cmd.AddCommand(
		newAuditListCmd(),
		newAuditGetCmd(),
		newAuditExportCmd(),
	)

	return cmd
}

func newAuditListCmd() *cobra.Command {
	var limit int
	var filter, from, to string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List audit events",
		RunE: func(cmd *cobra.Command, args []string) error {
			q := url.Values{}
			if limit > 0 {
				q.Set("$top", fmt.Sprintf("%d", limit))
			}
			if filter != "" {
				q.Set("$filter", filter)
			}
			if from != "" {
				q.Set("from", from)
			}
			if to != "" {
				q.Set("to", to)
			}
			raw, err := getClient().Get("/audit/events", q)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 0, "max events")
	cmd.Flags().StringVarP(&filter, "filter", "f", "", "filter expression")
	cmd.Flags().StringVar(&from, "from", "", "start date (RFC3339)")
	cmd.Flags().StringVar(&to, "to", "", "end date (RFC3339)")

	return cmd
}

func newAuditGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get an audit event by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := getClient().Get("/audit/events/"+args[0], nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newAuditExportCmd() *cobra.Command {
	var outputFile string

	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export audit trail",
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := getClient().Get("/audit/export", nil)
			if err != nil {
				return err
			}
			if outputFile != "" {
				return writeFile(outputFile, raw)
			}
			return printRaw(raw)
		},
	}

	cmd.Flags().StringVar(&outputFile, "output-file", "", "save to file")

	return cmd
}

// writeFile writes raw data to a file.
func writeFile(path string, data []byte) error {
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}
	fmt.Printf("Saved to %s\n", path)
	return nil
}
