package config

import (
	"os"
	"path/filepath"
	"runtime"
)

const (
	appName    = "action1"
	configFile = "config.yaml"
)

// Dir returns the configuration directory for the current OS.
func Dir() string {
	switch runtime.GOOS {
	case "darwin":
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "Library", "Application Support", appName)
	case "windows":
		return filepath.Join(os.Getenv("APPDATA"), appName)
	default: // linux, freebsd, etc.
		if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
			return filepath.Join(xdg, appName)
		}
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".config", appName)
	}
}

// FilePath returns the full path to the config file.
func FilePath() string {
	return filepath.Join(Dir(), configFile)
}
