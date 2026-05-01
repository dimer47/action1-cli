package cli

import (
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

func newSoftwareCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "software",
		Aliases: []string{"sw"},
		Short:   "Manage software repository",
	}

	cmd.AddCommand(
		newSWListCmd(),
		newSWCreateCmd(),
		newSWGetCmd(),
		newSWUpdateCmd(),
		newSWDeleteCmd(),
		newSWCloneCmd(),
		newSWMatchConflictsCmd(),
		newSWVersionCmd(),
		newSWUploadCmd(),
	)

	return cmd
}

func newSWListCmd() *cobra.Command {
	var limit int
	var filter string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List software repository packages",
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
			raw, err := getClient().Get("/software-repository/"+orgID, q)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 0, "max results")
	cmd.Flags().StringVarP(&filter, "filter", "f", "", "filter expression")

	return cmd
}

func newSWCreateCmd() *cobra.Command {
	var data string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new software package",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			body, err := parseDataFlag(data)
			if err != nil {
				return err
			}
			raw, err := getClient().Post("/software-repository/"+orgID, body)
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

func newSWGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <packageId>",
		Short: "Get software package settings",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Get(fmt.Sprintf("/software-repository/%s/%s", orgID, args[0]), nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newSWUpdateCmd() *cobra.Command {
	var data string

	cmd := &cobra.Command{
		Use:   "update <packageId>",
		Short: "Update a software package",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			body, err := parseDataFlag(data)
			if err != nil {
				return err
			}
			raw, err := getClient().Patch(fmt.Sprintf("/software-repository/%s/%s", orgID, args[0]), body)
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

func newSWDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <packageId>",
		Short: "Delete a custom software package",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			if !yes {
				if !confirmAction("delete software package " + args[0]) {
					return nil
				}
			}
			_, err := getClient().Delete(fmt.Sprintf("/software-repository/%s/%s", orgID, args[0]))
			if err != nil {
				return err
			}
			fmt.Println("Software package deleted.")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "skip confirmation")

	return cmd
}

func newSWCloneCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clone <packageId>",
		Short: "Clone a software package",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Post(fmt.Sprintf("/software-repository/%s/%s/clone", orgID, args[0]), nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newSWMatchConflictsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "match-conflicts [packageId]",
		Short: "Check for package matching conflicts",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			path := fmt.Sprintf("/software-repository/%s/match-conflicts", orgID)
			if len(args) > 0 {
				path = fmt.Sprintf("/software-repository/%s/%s/match-conflicts", orgID, args[0])
			}
			raw, err := getClient().Get(path, nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

// --- Version subcommands ---

func newSWVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Manage package versions",
	}

	cmd.AddCommand(
		newSWVersionCreateCmd(),
		newSWVersionGetCmd(),
		newSWVersionUpdateCmd(),
		newSWVersionDeleteCmd(),
		newSWVersionRemoveActionCmd(),
	)

	return cmd
}

func newSWVersionCreateCmd() *cobra.Command {
	var data string

	cmd := &cobra.Command{
		Use:   "create <packageId>",
		Short: "Create a new version",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			body, err := parseDataFlag(data)
			if err != nil {
				return err
			}
			raw, err := getClient().Post(fmt.Sprintf("/software-repository/%s/%s/versions", orgID, args[0]), body)
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

func newSWVersionGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <packageId> <versionId>",
		Short: "Get version details",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			raw, err := getClient().Get(fmt.Sprintf("/software-repository/%s/%s/versions/%s", orgID, args[0], args[1]), nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newSWVersionUpdateCmd() *cobra.Command {
	var data string

	cmd := &cobra.Command{
		Use:   "update <packageId> <versionId>",
		Short: "Update a version",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			body, err := parseDataFlag(data)
			if err != nil {
				return err
			}
			raw, err := getClient().Patch(fmt.Sprintf("/software-repository/%s/%s/versions/%s", orgID, args[0], args[1]), body)
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

func newSWVersionDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <packageId> <versionId>",
		Short: "Delete a version",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			if !yes {
				if !confirmAction("delete version " + args[1]) {
					return nil
				}
			}
			_, err := getClient().Delete(fmt.Sprintf("/software-repository/%s/%s/versions/%s", orgID, args[0], args[1]))
			if err != nil {
				return err
			}
			fmt.Println("Version deleted.")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "skip confirmation")

	return cmd
}

func newSWVersionRemoveActionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove-action <packageId> <versionId> <actionId>",
		Short: "Remove an additional action from a version",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}
			_, err := getClient().Delete(fmt.Sprintf("/software-repository/%s/%s/versions/%s/additional-actions/%s", orgID, args[0], args[1], args[2]))
			if err != nil {
				return err
			}
			fmt.Println("Action removed.")
			return nil
		},
	}
}

// --- Upload ---

func newSWUploadCmd() *cobra.Command {
	var chunkSize int

	cmd := &cobra.Command{
		Use:   "upload <packageId> <versionId> <filePath>",
		Short: "Upload an installation file (multi-chunk)",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := requireOrg(); err != nil {
				return err
			}

			packageID := args[0]
			versionID := args[1]
			filePath := args[2]

			file, err := os.Open(filePath)
			if err != nil {
				return fmt.Errorf("opening file: %w", err)
			}
			defer file.Close()

			stat, err := file.Stat()
			if err != nil {
				return fmt.Errorf("stat file: %w", err)
			}

			basePath := fmt.Sprintf("/software-repository/%s/%s/versions/%s/upload", orgID, packageID, versionID)

			// Initialize upload
			initBody := map[string]interface{}{
				"file_name": stat.Name(),
				"file_size": stat.Size(),
			}
			_, err = getClient().Post(basePath, initBody)
			if err != nil {
				return fmt.Errorf("initializing upload: %w", err)
			}

			// Upload chunks
			chunkBytes := chunkSize * 1024 * 1024
			buf := make([]byte, chunkBytes)
			offset := int64(0)

			for {
				n, readErr := file.Read(buf)
				if n > 0 {
					chunk := buf[:n]
					q := url.Values{
						"offset": {fmt.Sprintf("%d", offset)},
						"length": {fmt.Sprintf("%d", n)},
					}

					fullURL := basePath + "?" + q.Encode()
					_, err = getClient().Put(fullURL, io.NopCloser(io.NewSectionReader(file, offset, int64(n))), "application/octet-stream")
					if err != nil {
						return fmt.Errorf("uploading chunk at offset %d: %w", offset, err)
					}

					offset += int64(n)
					pct := float64(offset) / float64(stat.Size()) * 100
					fmt.Printf("\rUploading... %.1f%% (%d/%d bytes)", pct, offset, stat.Size())
					_ = chunk
				}
				if readErr == io.EOF {
					break
				}
				if readErr != nil {
					return fmt.Errorf("reading file: %w", readErr)
				}
			}

			fmt.Println("\nUpload complete.")
			return nil
		},
	}

	cmd.Flags().IntVar(&chunkSize, "chunk-size", 10, "chunk size in MB")

	return cmd
}
