package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newScriptCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "script",
		Short: "Manage script library",
	}

	cmd.AddCommand(
		newScriptListCmd(),
		newScriptCreateCmd(),
		newScriptGetCmd(),
		newScriptUpdateCmd(),
		newScriptDeleteCmd(),
	)

	return cmd
}

func newScriptListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List scripts",
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := getClient().Get("/scripts/all", nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newScriptCreateCmd() *cobra.Command {
	var data, name, description, scriptType, content, file string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a custom script",
		RunE: func(cmd *cobra.Command, args []string) error {
			var body map[string]interface{}

			if data != "" {
				var err error
				body, err = parseDataFlag(data)
				if err != nil {
					return err
				}
			} else {
				body = map[string]interface{}{}
			}

			if name != "" {
				body["name"] = name
			}
			if description != "" {
				body["description"] = description
			}
			if scriptType != "" {
				body["type"] = scriptType
			}
			if file != "" {
				fileContent, err := os.ReadFile(file)
				if err != nil {
					return fmt.Errorf("reading script file: %w", err)
				}
				body["content"] = string(fileContent)
			} else if content != "" {
				body["content"] = content
			}

			raw, err := getClient().Post("/scripts/all", body)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}

	cmd.Flags().StringVar(&data, "data", "", "JSON payload (overrides other flags)")
	cmd.Flags().StringVar(&name, "name", "", "script name")
	cmd.Flags().StringVar(&description, "description", "", "script description")
	cmd.Flags().StringVar(&scriptType, "type", "", "script type: powershell|cmd")
	cmd.Flags().StringVar(&content, "content", "", "script content (inline)")
	cmd.Flags().StringVarP(&file, "file", "f", "", "read script content from file")

	return cmd
}

func newScriptGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <scriptId>",
		Short: "Get a specific script",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := getClient().Get("/scripts/all/"+args[0], nil)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}
}

func newScriptUpdateCmd() *cobra.Command {
	var data, name, description, content, file string

	cmd := &cobra.Command{
		Use:   "update <scriptId>",
		Short: "Update a custom script",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var body map[string]interface{}

			if data != "" {
				var err error
				body, err = parseDataFlag(data)
				if err != nil {
					return err
				}
			} else {
				body = map[string]interface{}{}
			}

			if name != "" {
				body["name"] = name
			}
			if description != "" {
				body["description"] = description
			}
			if file != "" {
				fileContent, err := os.ReadFile(file)
				if err != nil {
					return fmt.Errorf("reading script file: %w", err)
				}
				body["content"] = string(fileContent)
			} else if content != "" {
				body["content"] = content
			}

			raw, err := getClient().Patch("/scripts/all/"+args[0], body)
			if err != nil {
				return err
			}
			return printRaw(raw)
		},
	}

	cmd.Flags().StringVar(&data, "data", "", "JSON payload")
	cmd.Flags().StringVar(&name, "name", "", "script name")
	cmd.Flags().StringVar(&description, "description", "", "script description")
	cmd.Flags().StringVar(&content, "content", "", "script content (inline)")
	cmd.Flags().StringVarP(&file, "file", "f", "", "read script content from file")

	return cmd
}

func newScriptDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <scriptId>",
		Short: "Delete a custom script",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes {
				if !confirmAction("delete script " + args[0]) {
					return nil
				}
			}
			_, err := getClient().Delete("/scripts/all/" + args[0])
			if err != nil {
				return err
			}
			fmt.Println("Script deleted.")
			return nil
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "skip confirmation")

	return cmd
}
