package auth

// Credentials holds the OAuth credentials and tokens.
type Credentials struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Store is the interface for credential storage backends.
type Store interface {
	Save(profile string, creds Credentials) error
	Load(profile string) (Credentials, error)
	Delete(profile string) error
}
