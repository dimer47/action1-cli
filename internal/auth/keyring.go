package auth

import (
	"encoding/json"
	"fmt"

	"github.com/zalando/go-keyring"

	"github.com/dimer47/action1-cli/internal/config"
)

const serviceName = "action1-cli"

// keyringSecret holds only the small secrets stored in the OS keychain.
type keyringSecret struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// KeyringStore stores client credentials in the OS keychain
// and tokens in the config file (tokens are too large for some keychains).
type KeyringStore struct {
	ConfigPath string
}

func NewKeyringStore(configPath string) *KeyringStore {
	return &KeyringStore{ConfigPath: configPath}
}

func (k *KeyringStore) Save(profile string, creds Credentials) error {
	// Store client_id + client_secret in keyring (small data)
	secret := keyringSecret{
		ClientID:     creds.ClientID,
		ClientSecret: creds.ClientSecret,
	}
	data, err := json.Marshal(secret)
	if err != nil {
		return fmt.Errorf("marshaling credentials: %w", err)
	}
	if err := keyring.Set(serviceName, profile, string(data)); err != nil {
		return fmt.Errorf("saving to keyring: %w", err)
	}

	// Store tokens in config file (can be large JWTs)
	cfg, err := config.Load(k.ConfigPath)
	if err != nil {
		return fmt.Errorf("loading config for token storage: %w", err)
	}
	p := cfg.Profiles[profile]
	p.AccessToken = creds.AccessToken
	p.RefreshToken = creds.RefreshToken
	cfg.Profiles[profile] = p
	if err := cfg.Save(k.ConfigPath); err != nil {
		return fmt.Errorf("saving tokens to config: %w", err)
	}

	return nil
}

func (k *KeyringStore) Load(profile string) (Credentials, error) {
	// Load client_id + client_secret from keyring
	data, err := keyring.Get(serviceName, profile)
	if err != nil {
		return Credentials{}, fmt.Errorf("loading credentials from keyring: %w", err)
	}

	var secret keyringSecret
	if err := json.Unmarshal([]byte(data), &secret); err != nil {
		return Credentials{}, fmt.Errorf("parsing credentials: %w", err)
	}

	// Load tokens from config file
	cfg, err := config.Load(k.ConfigPath)
	if err != nil {
		return Credentials{
			ClientID:     secret.ClientID,
			ClientSecret: secret.ClientSecret,
		}, nil
	}

	p, ok := cfg.Profiles[profile]
	if !ok {
		p = cfg.ActiveProfile()
	}

	return Credentials{
		ClientID:     secret.ClientID,
		ClientSecret: secret.ClientSecret,
		AccessToken:  p.AccessToken,
		RefreshToken: p.RefreshToken,
	}, nil
}

func (k *KeyringStore) Delete(profile string) error {
	// Remove from keyring
	_ = keyring.Delete(serviceName, profile)

	// Remove tokens from config
	cfg, err := config.Load(k.ConfigPath)
	if err != nil {
		return nil
	}
	p := cfg.Profiles[profile]
	p.AccessToken = ""
	p.RefreshToken = ""
	cfg.Profiles[profile] = p
	return cfg.Save(k.ConfigPath)
}
