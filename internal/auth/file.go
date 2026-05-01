package auth

import (
	"github.com/dimer47/action1-cli/internal/config"
)

// FileStore stores credentials in the config file (fallback when keyring is unavailable).
type FileStore struct {
	ConfigPath string
}

func NewFileStore(configPath string) *FileStore {
	return &FileStore{ConfigPath: configPath}
}

func (f *FileStore) Save(profile string, creds Credentials) error {
	cfg, err := config.Load(f.ConfigPath)
	if err != nil {
		return err
	}

	p := cfg.Profiles[profile]
	p.ClientID = creds.ClientID
	p.ClientSecret = creds.ClientSecret
	p.AccessToken = creds.AccessToken
	p.RefreshToken = creds.RefreshToken
	cfg.Profiles[profile] = p

	return cfg.Save(f.ConfigPath)
}

func (f *FileStore) Load(profile string) (Credentials, error) {
	cfg, err := config.Load(f.ConfigPath)
	if err != nil {
		return Credentials{}, err
	}

	p := cfg.ActiveProfile()
	if profile != "" {
		var ok bool
		p, ok = cfg.Profiles[profile]
		if !ok {
			return Credentials{}, nil
		}
	}

	return Credentials{
		ClientID:     p.ClientID,
		ClientSecret: p.ClientSecret,
		AccessToken:  p.AccessToken,
		RefreshToken: p.RefreshToken,
	}, nil
}

func (f *FileStore) Delete(profile string) error {
	cfg, err := config.Load(f.ConfigPath)
	if err != nil {
		return err
	}

	p := cfg.Profiles[profile]
	p.ClientID = ""
	p.ClientSecret = ""
	p.AccessToken = ""
	p.RefreshToken = ""
	cfg.Profiles[profile] = p

	return cfg.Save(f.ConfigPath)
}
