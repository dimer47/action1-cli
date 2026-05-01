package auth

import (
	"encoding/json"
	"fmt"

	"github.com/zalando/go-keyring"
)

// NewStore returns the best available credential store.
// It tries the OS keyring first, falling back to file-based storage.
func NewStore(configPath string, forceFile bool) Store {
	if forceFile {
		return NewFileStore(configPath)
	}

	// Try keyring by doing a test write/read/delete
	ks := NewKeyringStore(configPath)
	testKey := "__action1_cli_test__"
	testCreds := Credentials{ClientID: "test"}
	if err := keyringOnly(testKey, testCreds); err == nil {
		return ks
	}

	fmt.Println("Warning: OS keyring unavailable, using file-based credential storage")
	return NewFileStore(configPath)
}

// keyringOnly tests keyring access without touching the config file.
func keyringOnly(key string, creds Credentials) error {
	secret := keyringSecret{ClientID: creds.ClientID, ClientSecret: creds.ClientSecret}
	data, _ := json.Marshal(secret)
	if err := keyring.Set(serviceName, key, string(data)); err != nil {
		return err
	}
	_ = keyring.Delete(serviceName, key)
	return nil
}
