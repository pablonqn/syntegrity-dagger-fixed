package shared

import (
	"errors"
	"fmt"
	"os"
	"time"
)

// GitCredentials representa las credenciales para autenticación Git
type GitCredentials struct {
	User      string
	Token     string
	Source    string
	ExpiresAt time.Time
}

// CredentialSource indica de dónde vienen las credenciales
type CredentialSource string

const (
	SourceCI        CredentialSource = "CI"
	SourcePAT       CredentialSource = "PAT"
	SourceSSH       CredentialSource = "SSH"
	SourceAnonymous CredentialSource = "ANONYMOUS"
)

// ResolveGitCredentials intenta obtener credenciales de múltiples fuentes
func ResolveGitCredentials() (*GitCredentials, error) {
	// 1. Intentar con CI
	if os.Getenv("CI") == "true" {
		if token := os.Getenv("CI_JOB_TOKEN"); token != "" {
			return &GitCredentials{
				User:      "gitlab-ci-token",
				Token:     token,
				Source:    string(SourceCI),
				ExpiresAt: time.Now().Add(1 * time.Hour), // CI tokens suelen expirar en 1h
			}, nil
		}
	}

	// 2. Intentar con PAT
	if token := os.Getenv("GITLAB_PAT"); token != "" {
		return &GitCredentials{
			User:      "oauth2",
			Token:     token,
			Source:    string(SourcePAT),
			ExpiresAt: time.Now().Add(24 * time.Hour), // PATs suelen durar más
		}, nil
	}

	// 3. Intentar con SSH
	if key := os.Getenv("SSH_PRIVATE_KEY"); key != "" {
		return &GitCredentials{
			User:      "git",
			Token:     key,
			Source:    string(SourceSSH),
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}, nil
	}

	// 4. Modo anónimo (solo para repos públicos)
	fmt.Println("⚠️  No credentials found, using anonymous access")
	return &GitCredentials{
		User:      "",
		Token:     "",
		Source:    string(SourceAnonymous),
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}, nil
}

// Validate verifica si las credenciales son válidas.
func (c *GitCredentials) Validate() error {
	if c == nil {
		return errors.New("credentials are nil")
	}

	if c.ExpiresAt.Before(time.Now()) {
		return errors.New("credentials have expired")
	}

	switch CredentialSource(c.Source) {
	case SourceCI, SourcePAT:
		if c.User == "" || c.Token == "" {
			return errors.New("invalid credentials: user or token is empty")
		}
	case SourceSSH:
		if c.Token == "" {
			return errors.New("invalid SSH key")
		}
	case SourceAnonymous:
		// No validation needed
	default:
		return fmt.Errorf("unknown credential source: %s", c.Source)
	}

	return nil
}

// String implementa fmt.Stringer para logging seguro
func (c *GitCredentials) String() string {
	if c == nil {
		return "<nil>"
	}
	return fmt.Sprintf("GitCredentials{Source: %s, User: %s, ExpiresAt: %s}",
		c.Source, c.User, c.ExpiresAt.Format(time.RFC3339))
}
