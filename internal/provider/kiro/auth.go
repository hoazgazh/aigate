package kiro

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	kiroRefreshURLTemplate = "https://prod.%s.auth.desktop.kiro.dev/refreshToken"
	awsSSOOIDCURLTemplate  = "https://oidc.%s.amazonaws.com/token"
	kiroAPIHostTemplate    = "https://q.%s.amazonaws.com"
)

// AuthManager handles token lifecycle: load, refresh, cache.
type AuthManager struct {
	mu          sync.RWMutex
	creds       *Credentials
	fingerprint string
	client      *http.Client

	// Derived from region
	APIHost string
	QHost   string
}

// NewAuthManager creates an AuthManager from config.
// It auto-detects the credential source in priority order:
// SQLite DB > JSON file > env vars.
func NewAuthManager(refreshToken, profileARN, region, credsFile, cliDBFile string) (*AuthManager, error) {
	var creds *Credentials
	var err error

	switch {
	case cliDBFile != "":
		creds, err = LoadFromSQLite(cliDBFile)
		if err != nil {
			return nil, fmt.Errorf("load sqlite credentials: %w", err)
		}
		log.Printf("[auth] loaded credentials from SQLite (%s), type=%d", cliDBFile, creds.AuthType)

	case credsFile != "":
		creds, err = LoadFromJSON(credsFile)
		if err != nil {
			return nil, fmt.Errorf("credentials file error: %w\n\n  Check that the file exists and you are logged in to Kiro IDE or kiro-cli", err)
		}
		log.Printf("[auth] loaded credentials from JSON (%s), type=%d", credsFile, creds.AuthType)

	case refreshToken != "":
		creds = LoadFromEnv(refreshToken, profileARN, region)
		log.Printf("[auth] loaded credentials from env vars, type=%d", creds.AuthType)

	default:
		// Auto-detect: try common paths
		autoPath := autoDetectCredentials()
		if autoPath != "" {
			if strings.HasSuffix(autoPath, ".sqlite3") {
				creds, err = LoadFromSQLite(autoPath)
			} else {
				creds, err = LoadFromJSON(autoPath)
			}
			if err == nil {
				log.Printf("[auth] auto-detected credentials from %s, type=%d", autoPath, creds.AuthType)
			} else {
				return nil, fmt.Errorf("auto-detected %s but failed to load: %w", autoPath, err)
			}
		} else {
			return nil, fmt.Errorf("no credentials found")
		}
	}

	if creds.Region == "" {
		creds.Region = region
	}

	// Kiro API is only available in us-east-1 regardless of SSO region
	apiRegion := "us-east-1"

	am := &AuthManager{
		creds:       creds,
		fingerprint: machineFingerprint(),
		client:      &http.Client{Timeout: 30 * time.Second},
		APIHost:     fmt.Sprintf(kiroAPIHostTemplate, apiRegion),
		QHost:       fmt.Sprintf(kiroAPIHostTemplate, apiRegion),
	}

	log.Printf("[auth] sso_region=%s, api_region=%s, api_host=%s", creds.Region, apiRegion, am.APIHost)
	return am, nil
}

// GetToken returns a valid access token, refreshing if needed.
func (am *AuthManager) GetToken() (string, error) {
	am.mu.RLock()
	if am.creds.AccessToken != "" && !am.creds.IsExpiringSoon() {
		tok := am.creds.AccessToken
		am.mu.RUnlock()
		return tok, nil
	}
	am.mu.RUnlock()

	am.mu.Lock()
	defer am.mu.Unlock()

	// Double-check after acquiring write lock
	if am.creds.AccessToken != "" && !am.creds.IsExpiringSoon() {
		return am.creds.AccessToken, nil
	}

	// SQLite mode: reload first, kiro-cli might have refreshed
	if am.creds.sqliteDB != "" && am.creds.IsExpiringSoon() {
		reloaded, err := LoadFromSQLite(am.creds.sqliteDB)
		if err == nil && !reloaded.IsExpiringSoon() {
			am.creds.AccessToken = reloaded.AccessToken
			am.creds.RefreshToken = reloaded.RefreshToken
			am.creds.ExpiresAt = reloaded.ExpiresAt
			log.Printf("[auth] reloaded fresh token from SQLite")
			return am.creds.AccessToken, nil
		}
	}

	if err := am.refresh(); err != nil {
		// Graceful degradation: if refresh fails but token not yet expired, use it
		if am.creds.sqliteDB != "" && !am.creds.IsExpired() && am.creds.AccessToken != "" {
			log.Printf("[auth] refresh failed, using existing token until expiry: %v", err)
			return am.creds.AccessToken, nil
		}
		return "", fmt.Errorf("token refresh failed: %w", err)
	}

	return am.creds.AccessToken, nil
}

// ForceRefresh forces a token refresh (e.g. after 403).
func (am *AuthManager) ForceRefresh() (string, error) {
	am.mu.Lock()
	defer am.mu.Unlock()
	if err := am.refresh(); err != nil {
		return "", err
	}
	return am.creds.AccessToken, nil
}

// Headers returns the full set of Kiro API headers with a fresh token.
func (am *AuthManager) Headers() (map[string]string, error) {
	tok, err := am.GetToken()
	if err != nil {
		return nil, err
	}
	return kiroHeaders(tok, am.fingerprint), nil
}

// ProfileARN returns the profile ARN (may be empty for SSO).
func (am *AuthManager) ProfileARN() string {
	am.mu.RLock()
	defer am.mu.RUnlock()
	return am.creds.ProfileARN
}

// autoDetectCredentials tries common credential paths.
func autoDetectCredentials() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	paths := []string{
		filepath.Join(home, ".local", "share", "kiro-cli", "data.sqlite3"),
		filepath.Join(home, ".local", "share", "amazon-q", "data.sqlite3"),
		filepath.Join(home, ".aws", "sso", "cache", "kiro-auth-token.json"),
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

func (am *AuthManager) refresh() error {
	switch am.creds.AuthType {
	case AuthAWSSSOOIDC:
		return am.refreshAWSSSO()
	default:
		return am.refreshKiroDesktop()
	}
}

// refreshKiroDesktop refreshes via Kiro Desktop Auth endpoint.
// POST https://prod.{region}.auth.desktop.kiro.dev/refreshToken
// Body: {"refreshToken": "..."}
func (am *AuthManager) refreshKiroDesktop() error {
	if am.creds.RefreshToken == "" {
		return fmt.Errorf("refresh token is empty")
	}

	url := fmt.Sprintf(kiroRefreshURLTemplate, am.creds.Region)
	payload, _ := json.Marshal(map[string]string{
		"refreshToken": am.creds.RefreshToken,
	})

	req, _ := http.NewRequest("POST", url, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", fmt.Sprintf("KiroIDE-0.7.45-%s", am.fingerprint))

	resp, err := am.client.Do(req)
	if err != nil {
		return fmt.Errorf("kiro desktop refresh request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("kiro desktop refresh: status %d, body: %s", resp.StatusCode, body)
	}

	var result struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
		ExpiresIn    int    `json:"expiresIn"`
		ProfileARN   string `json:"profileArn"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("parse refresh response: %w", err)
	}
	if result.AccessToken == "" {
		return fmt.Errorf("no accessToken in response")
	}

	am.creds.AccessToken = result.AccessToken
	if result.RefreshToken != "" {
		am.creds.RefreshToken = result.RefreshToken
	}
	if result.ProfileARN != "" {
		am.creds.ProfileARN = result.ProfileARN
	}

	expiresIn := result.ExpiresIn
	if expiresIn == 0 {
		expiresIn = 3600
	}
	am.creds.ExpiresAt = time.Now().UTC().Add(time.Duration(expiresIn-60) * time.Second)

	log.Printf("[auth] token refreshed via Kiro Desktop, expires: %s", am.creds.ExpiresAt.Format(time.RFC3339))
	am.creds.Save()
	return nil
}

// refreshAWSSSO refreshes via AWS SSO OIDC endpoint.
// POST https://oidc.{region}.amazonaws.com/token
// Body (JSON): {"grantType":"refresh_token","clientId":"...","clientSecret":"...","refreshToken":"..."}
func (am *AuthManager) refreshAWSSSO() error {
	err := am.doAWSSSORefresh()
	if err == nil {
		return nil
	}

	// On 400 (stale token), reload from SQLite and retry once
	if am.creds.sqliteDB != "" {
		log.Printf("[auth] SSO refresh failed, reloading from SQLite and retrying: %v", err)
		reloaded, loadErr := LoadFromSQLite(am.creds.sqliteDB)
		if loadErr == nil {
			am.creds.RefreshToken = reloaded.RefreshToken
			am.creds.ClientID = reloaded.ClientID
			am.creds.ClientSecret = reloaded.ClientSecret
			return am.doAWSSSORefresh()
		}
	}
	return err
}

func (am *AuthManager) doAWSSSORefresh() error {
	if am.creds.RefreshToken == "" {
		return fmt.Errorf("refresh token is empty")
	}
	if am.creds.ClientID == "" || am.creds.ClientSecret == "" {
		return fmt.Errorf("client_id/client_secret required for AWS SSO OIDC")
	}

	region := am.creds.SSORegion
	if region == "" {
		region = am.creds.Region
	}
	url := fmt.Sprintf(awsSSOOIDCURLTemplate, region)

	payload, _ := json.Marshal(map[string]string{
		"grantType":    "refresh_token",
		"clientId":     am.creds.ClientID,
		"clientSecret": am.creds.ClientSecret,
		"refreshToken": am.creds.RefreshToken,
	})

	req, _ := http.NewRequest("POST", url, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := am.client.Do(req)
	if err != nil {
		return fmt.Errorf("aws sso refresh request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("aws sso refresh: status %d, body: %s", resp.StatusCode, body)
	}

	var result struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
		ExpiresIn    int    `json:"expiresIn"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("parse sso response: %w", err)
	}
	if result.AccessToken == "" {
		return fmt.Errorf("no accessToken in SSO response")
	}

	am.creds.AccessToken = result.AccessToken
	if result.RefreshToken != "" {
		am.creds.RefreshToken = result.RefreshToken
	}

	expiresIn := result.ExpiresIn
	if expiresIn == 0 {
		expiresIn = 3600
	}
	am.creds.ExpiresAt = time.Now().UTC().Add(time.Duration(expiresIn-60) * time.Second)

	log.Printf("[auth] token refreshed via AWS SSO OIDC, expires: %s", am.creds.ExpiresAt.Format(time.RFC3339))
	am.creds.Save()
	return nil
}
