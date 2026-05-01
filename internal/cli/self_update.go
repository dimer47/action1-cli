package cli

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/cobra"

	"github.com/dimer47/action1-cli/internal/update"
)

func newSelfUpdateCmd() *cobra.Command {
	var check bool

	cmd := &cobra.Command{
		Use:   "self-update",
		Short: "Update the CLI to the latest version",
		RunE: func(cmd *cobra.Command, args []string) error {
			if Version == "dev" {
				return fmt.Errorf("self-update is not available for dev builds — use 'go build' instead")
			}

			result := update.Check(Version)

			if check {
				if result != nil && result.HasUpdate {
					fmt.Printf("New version available: %s → %s\n", result.CurrentVersion, result.LatestVersion)
				} else {
					fmt.Printf("Already up to date (v%s)\n", Version)
				}
				return nil
			}

			if result == nil || !result.HasUpdate {
				fmt.Printf("Already up to date (v%s)\n", Version)
				return nil
			}

			fmt.Printf("Updating action1: %s → %s\n", result.CurrentVersion, result.LatestVersion)

			// Determine download URL
			goos := runtime.GOOS
			goarch := runtime.GOARCH
			ext := "tar.gz"
			if goos == "windows" {
				ext = "zip"
			}
			url := fmt.Sprintf(
				"https://github.com/dimer47/action1-cli/releases/latest/download/action1-cli_%s_%s.%s",
				goos, goarch, ext,
			)

			// Download to temp file
			fmt.Printf("Downloading from %s...\n", url)
			client := &http.Client{Timeout: 120 * time.Second}
			resp, err := client.Get(url)
			if err != nil {
				return fmt.Errorf("download failed: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("download failed: HTTP %d", resp.StatusCode)
			}

			tmpDir, err := os.MkdirTemp("", "action1-update-*")
			if err != nil {
				return fmt.Errorf("creating temp dir: %w", err)
			}
			defer os.RemoveAll(tmpDir)

			archivePath := filepath.Join(tmpDir, "action1."+ext)
			f, err := os.Create(archivePath)
			if err != nil {
				return fmt.Errorf("creating temp file: %w", err)
			}

			written, err := io.Copy(f, resp.Body)
			f.Close()
			if err != nil {
				return fmt.Errorf("downloading: %w", err)
			}
			fmt.Printf("Downloaded %.1f MB\n", float64(written)/1024/1024)

			// Extract
			binaryName := "action1"
			if goos == "windows" {
				binaryName = "action1.exe"
			}

			if ext == "tar.gz" {
				extractCmd := exec.Command("tar", "xzf", archivePath, "-C", tmpDir)
				if out, err := extractCmd.CombinedOutput(); err != nil {
					return fmt.Errorf("extracting archive: %s: %w", string(out), err)
				}
			} else {
				extractCmd := exec.Command("unzip", "-o", archivePath, "-d", tmpDir)
				if out, err := extractCmd.CombinedOutput(); err != nil {
					return fmt.Errorf("extracting archive: %s: %w", string(out), err)
				}
			}

			newBinary := filepath.Join(tmpDir, binaryName)
			if _, err := os.Stat(newBinary); os.IsNotExist(err) {
				return fmt.Errorf("binary not found in archive")
			}

			// Find current binary location
			currentBinary, err := os.Executable()
			if err != nil {
				return fmt.Errorf("finding current binary: %w", err)
			}
			currentBinary, err = filepath.EvalSymlinks(currentBinary)
			if err != nil {
				return fmt.Errorf("resolving symlinks: %w", err)
			}

			// Try direct replace first
			fmt.Printf("Installing to %s...\n", currentBinary)
			if err := replaceBinary(newBinary, currentBinary); err != nil {
				// Need sudo
				fmt.Println("Permission denied — trying with sudo...")
				sudoCmd := exec.Command("sudo", "cp", newBinary, currentBinary)
				sudoCmd.Stdin = os.Stdin
				sudoCmd.Stdout = os.Stdout
				sudoCmd.Stderr = os.Stderr
				if err := sudoCmd.Run(); err != nil {
					return fmt.Errorf("install failed: %w\nYou can manually run: sudo cp %s %s", err, newBinary, currentBinary)
				}
			}

			// Clear update cache
			update.ClearCache()

			fmt.Printf("Successfully updated to v%s\n", result.LatestVersion)
			return nil
		},
	}

	cmd.Flags().BoolVar(&check, "check", false, "only check for updates without installing")

	return cmd
}

func replaceBinary(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return dstFile.Chmod(0755)
}
