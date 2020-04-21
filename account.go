package eapi

import (
	"net/url"
	"strings"
	"time"
)

const (
	// Authorization (basic) tokens
	eglToken = "MzRhMDJjZjhmNDQxNGUyOWIxNTkyMTg3NmRhMzZmOWE6ZGFhZmJjY2M3Mzc3NDUwMzlkZmZlNTNkOTRmYzc2Y2Y="
	fnToken  = "ZWM2ODRiOGM2ODdmNDc5ZmFkZWEzY2IyYWQ4M2Y1YzY6ZTFmMzFjMjExZjI4NDEzMTg2MjYyZDM3YTEzZmM4NGQ="

	deviceID = "70136f376c6979de6ead69be74ef6542"
)

// AccountService handles all account-related routes.
type AccountService service

// TokenResponse represents an OAuth2 token response.
type TokenResponse struct {
	AccessToken    string    `json:"access_token"`
	ClientID       string    `json:"client_id"`
	ClientService  string    `json:"client_service"`
	ExpiresIn      int       `json:"expires_in"`
	ExpiresAt      time.Time `json:"expires_at"`
	TokenType      string    `json:"token_type"`
	InternalClient bool      `json:"internal_client"`
}

// ExchangeTokenResponse represents an extended OAuth2 token response.
type ExchangeTokenResponse struct {
	TokenResponse
	AccountID        string    `json:"account_id"`
	App              string    `json:"app"`
	InAppID          string    `json:"in_app_id"`
	DeviceID         string    `json:"device_id"`
	DisplayName      string    `json:"displayName"`
	RefreshToken     string    `json:"refresh_token"`
	RefreshExpires   int       `json:"refresh_expires"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
}

// AuthWithClientCredentials authenticates the client with client credentials.
func (s *AccountService) AuthWithClientCredentials() (res *TokenResponse, err error) {
	// Build request body
	body := &url.Values{}
	body.Add("grant_type", "client_credentials")
	body.Add("token_type", "eg1")

	// Create request
	req, err := s.client.newReq("POST", accountURL+"/oauth/token", strings.NewReader(body.Encode()))
	if err != nil {
		return
	}

	// Set custom headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if s.client.opts.Type == EGL {
		req.Header.Set("Authorization", "basic "+eglToken)
	} else if s.client.opts.Type == FN {
		req.Header.Set("Authorization", "basic "+fnToken)
	}

	// Make request
	res = &TokenResponse{}
	err = s.client.do(req, res)

	// Set session info
	if err == nil {
		s.client.session.token = res.AccessToken
		s.client.session.tokenExpires = res.ExpiresAt
	}

	return
}

// AuthWithExchangeCode authenticates the client with a provided exchange code.
func (s *AccountService) AuthWithExchangeCode(exchangeCode string) (res *ExchangeTokenResponse, err error) {
	// Build request body
	body := &url.Values{}
	body.Add("grant_type", "exchange_code")
	body.Add("token_type", "eg1")
	body.Add("exchange_code", exchangeCode)

	// Create request
	req, err := s.client.newReq("POST", accountURL+"/oauth/token", strings.NewReader(body.Encode()))
	if err != nil {
		return
	}

	// Set custom headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Epic-Device-ID", deviceID)
	if s.client.opts.Type == EGL {
		req.Header.Set("Authorization", "basic "+eglToken)
	} else if s.client.opts.Type == FN {
		req.Header.Set("Authorization", "basic "+fnToken)
	}

	// Make request
	res = &ExchangeTokenResponse{}
	err = s.client.do(req, res)

	// Set session info
	if err == nil {
		s.client.session.token = res.AccessToken
		s.client.session.tokenExpires = res.ExpiresAt

		s.client.session.refreshToken = res.RefreshToken
		s.client.session.refreshExpires = res.RefreshExpiresAt

		s.client.session.accountID = res.AccountID
		s.client.session.deviceID = res.DeviceID
	}

	return
}
