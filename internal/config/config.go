package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Region represents a server region.
type Region string

const (
	RegionNA Region = "na"
	RegionEU Region = "eu"
	RegionAU Region = "au"
)

// BaseURL returns the API base URL for the region.
func (r Region) BaseURL() string {
	switch r {
	case RegionEU:
		return "https://app.eu.action1.com/api/3.0"
	case RegionAU:
		return "https://app.au.action1.com/api/3.0"
	default:
		return "https://app.action1.com/api/3.0"
	}
}

// Profile holds per-profile configuration.
type Profile struct {
	Region Region `yaml:"region,omitempty"`
	Org    string `yaml:"org,omitempty"`
	Output string `yaml:"output,omitempty"`

	// Fallback token storage (when keyring is unavailable).
	AccessToken  string `yaml:"access_token,omitempty"`
	RefreshToken string `yaml:"refresh_token,omitempty"`
	ClientID     string `yaml:"client_id,omitempty"`
	ClientSecret string `yaml:"client_secret,omitempty"`
}

// Config is the root configuration structure.
type Config struct {
	CurrentProfile string             `yaml:"current_profile"`
	Profiles       map[string]Profile `yaml:"profiles"`
}

// Load reads the config from disk. Returns a default config if the file doesn't exist.
func Load(path string) (*Config, error) {
	if path == "" {
		path = FilePath()
	}

	cfg := &Config{
		CurrentProfile: "default",
		Profiles: map[string]Profile{
			"default": {Region: RegionNA, Output: "table"},
		},
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if cfg.Profiles == nil {
		cfg.Profiles = map[string]Profile{
			"default": {Region: RegionNA, Output: "table"},
		}
	}

	return cfg, nil
}

// Save writes the config to disk.
func (c *Config) Save(path string) error {
	if path == "" {
		path = FilePath()
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	return nil
}

// ActiveProfile returns the currently active profile.
func (c *Config) ActiveProfile() Profile {
	p, ok := c.Profiles[c.CurrentProfile]
	if !ok {
		return Profile{Region: RegionNA, Output: "table"}
	}
	return p
}

// SetProfileValue sets a key-value pair on a profile.
func (c *Config) SetProfileValue(profile, key, value string) error {
	p, ok := c.Profiles[profile]
	if !ok {
		p = Profile{}
	}

	switch key {
	case "region":
		p.Region = Region(value)
	case "org":
		p.Org = value
	case "output":
		p.Output = value
	default:
		return fmt.Errorf("unknown config key: %s", key)
	}

	c.Profiles[profile] = p
	return nil
}

// GetProfileValue returns the value for a config key.
func (c *Config) GetProfileValue(profile, key string) (string, error) {
	p, ok := c.Profiles[profile]
	if !ok {
		return "", fmt.Errorf("profile %q not found", profile)
	}

	switch key {
	case "region":
		return string(p.Region), nil
	case "org":
		return p.Org, nil
	case "output":
		return p.Output, nil
	default:
		return "", fmt.Errorf("unknown config key: %s", key)
	}
}
