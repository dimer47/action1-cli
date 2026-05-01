package cli

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

func newReportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report",
		Short: "Manage reports and report data",
	}

	cmd.AddCommand(
		newReportListCmd(),
		newReportCreateCmd(),
		newReportUpdateCmd(),
		newReportDeleteCmd(),
		newReportDataCmd(),
		newReportErrorsCmd(),
		newReportExportCmd(),
		newReportRequeryCmd(),
		newReportDrilldownCmd(),
		newReportDrilldownExportCmd(),
	)

	return cmd
}

func newReportListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list [categoryId]",
		Short: "List reports (optionally by category)",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "/reports/all"
			if len(args) > 0 {
				path = "/reports/all/" + args[0]
			}
			raw, err := getClient().Get(path, nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newReportCreateCmd() *cobra.Command {
	var data string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a custom report",
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := parseDataFlag(data)
			if err != nil {
				return err
			}
			raw, err := getClient().Post("/reports/all/custom", body)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}

	cmd.Flags().StringVar(&data, "data", "", "JSON payload")
	_ = cmd.MarkFlagRequired("data")

	return cmd
}

func newReportUpdateCmd() *cobra.Command {
	var data string

	cmd := &cobra.Command{
		Use:   "update <reportId>",
		Short: "Update a custom report",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := parseDataFlag(data)
			if err != nil {
				return err
			}
			raw, err := getClient().Patch("/reports/all/custom/"+args[0], body)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}

	cmd.Flags().StringVar(&data, "data", "", "JSON payload")
	_ = cmd.MarkFlagRequired("data")

	return cmd
}

func newReportDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <reportId>",
		Short: "Delete a custom report",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes {
				if !confirmAction("delete report " + args[0]) {
					return nil
				}
			}
			_, err := getClient().Delete("/reports/all/custom/" + args[0])
			if err != nil {
				return err
			}
			fmt.Println("Report deleted.")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "skip confirmation")

	return cmd
}

func newReportDataCmd() *cobra.Command {
	var limit int
	var filter string
	var all bool

	cmd := &cobra.Command{
		Use:   "data <reportId>",
		Short: "Get report rows",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			q := url.Values{}
			if limit > 0 {
				q.Set("$top", fmt.Sprintf("%d", limit))
			}
			if filter != "" {
				q.Set("$filter", filter)
			}

			path := fmt.Sprintf("/reportdata/%s/%s/data", orgID, args[0])
			if all {
				items, err := getClient().GetAll(path, q)
				if err != nil {
					return err
				}
				return printResult(rawToInterface(items))
			}

			raw, err := getClient().Get(path, q)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 0, "max rows")
	cmd.Flags().StringVarP(&filter, "filter", "f", "", "filter expression")
	cmd.Flags().BoolVar(&all, "all", false, "fetch all rows")

	return cmd
}

func newReportErrorsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "errors <reportId>",
		Short: "Get report errors",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Get(fmt.Sprintf("/reportdata/%s/%s/errors", orgID, args[0]), nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newReportExportCmd() *cobra.Command {
	var outputFile string

	cmd := &cobra.Command{
		Use:   "export <reportId>",
		Short: "Export a report",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Get(fmt.Sprintf("/reportdata/%s/%s/export", orgID, args[0]), nil)
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

func newReportRequeryCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "requery <reportId>",
		Short: "Re-query a report",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Post(fmt.Sprintf("/reportdata/%s/%s/requery", orgID, args[0]), nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newReportDrilldownCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "drilldown <reportId> <rowId>",
		Short: "Drill down to report row details",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Get(fmt.Sprintf("/reportdata/%s/%s/data/%s/drilldown", orgID, args[0], args[1]), nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newReportDrilldownExportCmd() *cobra.Command {
	var outputFile string

	cmd := &cobra.Command{
		Use:   "drilldown-export <reportId> <rowId>",
		Short: "Export drilldown details",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Get(fmt.Sprintf("/reportdata/%s/%s/data/%s/export", orgID, args[0], args[1]), nil)
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
