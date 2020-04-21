package eapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// User agents
	eglUA = "UELauncher/10.14.0-11896852+++Portal+Release-Live Windows/10.0.18363.1.256.64bit"
	fnUA  = "Fortnite/++Fortnite+Release-12.20-CL-12170032 Windows/10.0.18363.1.256.64bit"

	// URLs
	accountURL  = "https://account-public-service-prod.ol.epicgames.com/account/api"
	fortniteURL = "https://fortnite-public-service-prod11.ol.epicgames.com/fortnite/api"
)

// ClientOptions defines user options.
type ClientOptions struct {
	UpdateTokens bool
	Type         ClientType
}

// Client implements the API.
type Client struct {
	opts *ClientOptions

	httpClient *http.Client

	session authSession

	common service

	// Services
	Account  *AccountService
	Fortnite *FortniteService
}

type service struct {
	client *Client
}

type authSession struct {
	token        string
	refreshToken string

	tokenExpires   time.Time
	refreshExpires time.Time

	accountID string
	deviceID  string
}

// ClientType represents the client type
type ClientType uint8

const (
	// EGL (Epic Games Launcher)
	EGL ClientType = iota
	// FN (Fortnite)
	FN
)

type errorResponse struct {
	Code    string `json:"errorCode"`
	Message string `json:"message"`
}

// NewClient returns a new API instance.
func NewClient(options *ClientOptions) *Client {
	if options == nil {
		options = &ClientOptions{}
	}

	c := &Client{opts: options, httpClient: &http.Client{Timeout: time.Duration(30 * time.Second)}}
	c.common.client = c

	c.Account = (*AccountService)(&c.common)
	c.Fortnite = (*FortniteService)(&c.common)

	// Keep tokens up-to-date
	if options.UpdateTokens {
		go func() {
			time.Sleep(60 * time.Second)
		}()
	}

	return c
}

func (c *Client) newReq(method string, url string, body io.Reader) (req *http.Request, err error) {
	// Create HTTP request
	req, err = http.NewRequest(method, url, body)
	if err != nil {
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	if c.opts.Type == EGL {
		req.Header.Set("User-Agent", eglUA)
	} else if c.opts.Type == FN {
		req.Header.Set("User-Agent", fnUA)
	}
	if c.session.token != "" {
		req.Header.Set("Authorization", "bearer "+c.session.token)
	}

	return
}

func (c *Client) do(req *http.Request, res interface{}) error {
	// Make HTTP request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		if resp.Header.Get("Content-Type") == "application/json" {
			respError := &errorResponse{}
			if err := json.NewDecoder(resp.Body).Decode(respError); err != nil {
				return fmt.Errorf("failed to decode error body: %v", err)
			}

			return fmt.Errorf("%s: %s", respError.Code, respError.Message)
		}
		return fmt.Errorf("bad status-code: %d", resp.StatusCode)
	}

	// Parse response
	if resp.StatusCode == http.StatusOK && resp.Header.Get("Content-Type") == "application/json" && res != nil {
		if err := json.NewDecoder(resp.Body).Decode(res); err != nil && err != io.EOF {
			return err
		}
	}

	return nil
}
