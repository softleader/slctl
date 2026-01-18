package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/go-github/v69/github"
)

const (
	deviceCodeURL   = "https://github.com/login/device/code"
	accessTokenURL  = "https://github.com/login/oauth/access_token"
	defaultClientID = "Ov23li1SPkIjssKOzz2f"
)

var (
	// Scopes represents the scopes required for slctl
	Scopes = []github.Scope{github.ScopeReadOrg, github.ScopeUser, github.ScopeRepo}
)

// DeviceCodeResponse represents the response from GitHub device code endpoint
type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

// AccessTokenResponse represents the response from GitHub access token endpoint
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	Error       string `json:"error"`
}

// RequestDeviceCode requests a device code from GitHub
func RequestDeviceCode(ctx context.Context, clientID string, scopes []github.Scope) (*DeviceCodeResponse, error) {
	if clientID == "" {
		clientID = defaultClientID
	}
	if clientID == "" {
		return nil, fmt.Errorf("client ID is required for Device Flow")
	}

	form := url.Values{}
	form.Add("client_id", clientID)
	var scopeStrings []string
	for _, s := range scopes {
		scopeStrings = append(scopeStrings, string(s))
	}
	form.Add("scope", strings.Join(scopeStrings, " "))

	req, err := http.NewRequestWithContext(ctx, "POST", deviceCodeURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var dcr DeviceCodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&dcr); err != nil {
		return nil, err
	}

	return &dcr, nil
}

// PollAccessToken polls GitHub for an access token
func PollAccessToken(ctx context.Context, clientID, deviceCode string, interval int) (string, error) {
	if clientID == "" {
		clientID = defaultClientID
	}
	if interval <= 0 {
		interval = 5
	}

	form := url.Values{}
	form.Add("client_id", clientID)
	form.Add("device_code", deviceCode)
	form.Add("grant_type", "urn:ietf:params:oauth:grant-type:device_code")

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-ticker.C:
			req, err := http.NewRequestWithContext(ctx, "POST", accessTokenURL, strings.NewReader(form.Encode()))
			if err != nil {
				return "", err
			}
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return "", err
			}

			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				return "", err
			}

			var atr AccessTokenResponse
			if err := json.Unmarshal(body, &atr); err != nil {
				return "", err
			}

			if atr.AccessToken != "" {
				return atr.AccessToken, nil
			}

			switch atr.Error {
			case "authorization_pending":
				// Continue polling
			case "slow_down":
				interval += 5
				ticker.Reset(time.Duration(interval) * time.Second)
			case "expired_token":
				return "", fmt.Errorf("the device code has expired")
			case "access_denied":
				return "", fmt.Errorf("access denied by user")
			default:
				return "", fmt.Errorf("oauth error: %s", atr.Error)
			}
		}
	}
}
