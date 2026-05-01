package cli

import (
	"net/url"

	"github.com/spf13/cobra"
)

func newSearchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "search <query>",
		Short: "Search reports, endpoints, and apps",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			q := url.Values{
				"q": {args[0]},
			}
			raw, err := getClient().Get("/search/"+orgID, q)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}
