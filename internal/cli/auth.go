package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/dimer47/action1-cli/internal/auth"
	"github.com/dimer47/action1-cli/internal/config"
)

func newAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication",
	}

	cmd.AddCommand(
		newAuthLoginCmd(),
		newAuthLogoutCmd(),
		newAuthStatusCmd(),
		newAuthTokenCmd(),
		newAuthRefreshCmd(),
	)

	return cmd
}

func newAuthLoginCmd() *cobra.Command {
	var clientID, clientSecret string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with Action1 API",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Region: %s (%s)\n", region, config.Region(region).BaseURL())

			// Interactive mode if not provided via flags
			if clientID == "" {
				fmt.Print("Client ID: ")
				scanner := bufio.NewScanner(os.Stdin)
				if scanner.Scan() {
					clientID = strings.TrimSpace(scanner.Text())
				}
			}
			if clientSecret == "" {
				fmt.Print("Client Secret: ")
				secretBytes, err := term.ReadPassword(int(syscall.Stdin))
				if err != nil {
					return fmt.Errorf("reading secret: %w", err)
				}
				fmt.Println()
				clientSecret = string(secretBytes)
			}

			if clientID == "" || clientSecret == "" {
				return fmt.Errorf("client ID and secret are required")
			}

			client := getClient()
			oauthResp, err := client.Authenticate(clientID, clientSecret)
			if err != nil {
				return fmt.Errorf("authentication failed: %w", err)
			}

			creds := auth.Credentials{
				ClientID:     clientID,
				ClientSecret: clientSecret,
				AccessToken:  oauthResp.AccessToken,
				RefreshToken: oauthResp.RefreshToken,
			}

			if err := store.Save(profile, creds); err != nil {
				return fmt.Errorf("saving credentials: %w", err)
			}

			// Persist the region in config so subsequent commands use it
			if err := cfg.SetProfileValue(cfg.CurrentProfile, "region", region); err == nil {
				_ = cfg.Save(cfgFile)
			}

			fmt.Println("Successfully authenticated.")
			fmt.Printf("Token expires in %d seconds.\n", oauthResp.ExpiresIn)
			return nil
		},
	}

	cmd.Flags().StringVar(&clientID, "client-id", "", "API Client ID")
	cmd.Flags().StringVar(&clientSecret, "client-secret", "", "API Client Secret")

	return cmd
}

func newAuthLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Remove stored credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := store.Delete(profile); err != nil {
				return fmt.Errorf("removing credentials: %w", err)
			}
			fmt.Println("Credentials removed.")
			return nil
		},
	}
}

func newAuthStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show authentication status",
		RunE: func(cmd *cobra.Command, args []string) error {
			creds, err := store.Load(profile)
			if err != nil || creds.AccessToken == "" {
				fmt.Println("Not authenticated.")
				fmt.Println("Run 'action1 auth login' to authenticate.")
				return nil
			}

			fmt.Printf("Profile:   %s\n", profile)
			fmt.Printf("Region:    %s\n", region)
			fmt.Printf("Client ID: %s...%s\n", creds.ClientID[:4], creds.ClientID[len(creds.ClientID)-4:])
			fmt.Println("Status:    Authenticated")

			if creds.RefreshToken != "" {
				fmt.Println("Refresh:   Available")
			}

			return nil
		},
	}
}

func newAuthTokenCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "token",
		Short: "Print the current access token",
		RunE: func(cmd *cobra.Command, args []string) error {
			creds, err := store.Load(profile)
			if err != nil || creds.AccessToken == "" {
				return fmt.Errorf("not authenticated — run 'action1 auth login' first")
			}
			fmt.Print(creds.AccessToken)
			return nil
		},
	}
}

func newAuthRefreshCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "refresh",
		Short: "Force token refresh",
		RunE: func(cmd *cobra.Command, args []string) error {
			creds, err := store.Load(profile)
			if err != nil {
				return fmt.Errorf("not authenticated — run 'action1 auth login' first")
			}

			if creds.RefreshToken == "" {
				return fmt.Errorf("no refresh token available — run 'action1 auth login' again")
			}

			client := getClient()
			oauthResp, err := client.RefreshAuth(creds.ClientID, creds.ClientSecret, creds.RefreshToken)
			if err != nil {
				return fmt.Errorf("token refresh failed: %w", err)
			}

			creds.AccessToken = oauthResp.AccessToken
			creds.RefreshToken = oauthResp.RefreshToken

			if err := store.Save(profile, creds); err != nil {
				return fmt.Errorf("saving credentials: %w", err)
			}

			fmt.Println("Token refreshed successfully.")
			fmt.Printf("Expires in %d seconds.\n", oauthResp.ExpiresIn)
			return nil
		},
	}
}
