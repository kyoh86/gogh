package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/cli/browser"
	"github.com/kyoh86/gogh/v2"
	"github.com/kyoh86/gogh/v2/internal/github"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var loginFlags struct {
	Host string
}

const clientID = "Ov23li6aEWIxek6F8P5L"

type DeviceCodeResponse struct {
	DeviceCode              string
	UserCode                string
	VerificationURI         string
	VerificationURIComplete string
	ExpiresIn               int
	Interval                int
}

type TokenResponse struct {
	AccessToken string
	Scope       string
	TokenType   string
}

type ErrorResponse struct {
	Error            string
	ErrorDescription string
	ErrorURI         string
}

var scopes = []string{"repo", "delete_repo"}

var loginCommand = &cobra.Command{
	Use:     "login",
	Aliases: []string{"signin", "add"},
	Short:   "Login for the host and owner",
	Args:    cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		if loginFlags.Host == "" {
			loginFlags.Host = github.DefaultHost
			if err := huh.NewForm(huh.NewGroup(
				huh.NewInput().
					Title("Host name").
					Validate(gogh.ValidateHost).
					Value(&loginFlags.Host),
			)).Run(); err != nil {
				return err
			}
		}
		host := loginFlags.Host

		oauthConfig := &oauth2.Config{
			ClientID: clientID,
			Endpoint: oauth2.Endpoint{
				AuthURL:  fmt.Sprintf("https://%s/login/device/code", host),
				TokenURL: fmt.Sprintf("https://%s/login/oauth/access_token", host),
			},
			Scopes: scopes,
		}

		// Request device code
		deviceCodeResp, err := requestDeviceCode(oauthConfig.ClientID, oauthConfig.Endpoint.AuthURL)
		if err != nil {
			return fmt.Errorf("failed to request device code: %w", err)
		}

		if errors.Is(browser.OpenURL(deviceCodeResp.VerificationURI), exec.ErrNotFound) {
			fmt.Printf("Visit %s and enter the code: %s\n", deviceCodeResp.VerificationURI, deviceCodeResp.UserCode)
		} else {
			fmt.Printf("Opened %s, so enter the code: %s\n", deviceCodeResp.VerificationURI, deviceCodeResp.UserCode)
		}

		// Poll for token
		tokenResp, err := pollForToken(oauthConfig, deviceCodeResp)
		if err != nil {
			return fmt.Errorf("failed to poll for token: %w", err)
		}

		// Get user info
		adaptor, err := github.NewAdaptor(context.Background(), host, tokenResp.AccessToken)
		if err != nil {
			return fmt.Errorf("failed to create GitHub adaptor: %w", err)
		}
		user, _, err := adaptor.GetAuthenticatedUser(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get authenticated user info: %w", err)
		}

		tokens.Set(host, user.GetLogin(), tokenResp.AccessToken)

		fmt.Println("Login successful!")
		return nil
	},
}

func requestDeviceCode(clientID, authURL string) (*DeviceCodeResponse, error) {
	resp, err := http.PostForm(authURL, url.Values{
		"client_id": {clientID},
		"scope":     scopes,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, err
	}

	deviceCodeResp := &DeviceCodeResponse{
		DeviceCode:              values.Get("device_code"),
		UserCode:                values.Get("user_code"),
		VerificationURI:         values.Get("verification_uri"),
		VerificationURIComplete: values.Get("verification_uri_complete"),
		ExpiresIn:               atoi(values.Get("expires_in")),
		Interval:                atoi(values.Get("interval")),
	}

	return deviceCodeResp, nil
}

func atoi(str string) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return val
}

func pollForToken(oauthConfig *oauth2.Config, deviceCodeResp *DeviceCodeResponse) (*TokenResponse, error) {
POLL_LOOP:
	for {
		time.Sleep(time.Duration(deviceCodeResp.Interval*2) * time.Second) // Intervalを2倍にしてリクエスト頻度を下げる

		resp, err := http.PostForm(oauthConfig.Endpoint.TokenURL, url.Values{
			"client_id":   {oauthConfig.ClientID},
			"device_code": {deviceCodeResp.DeviceCode},
			"grant_type":  {"urn:ietf:params:oauth:grant-type:device_code"},
		})
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		values, err := url.ParseQuery(string(body))
		if err != nil {
			return nil, err
		}

		switch {
		case resp.StatusCode == http.StatusOK && values.Get("error") == "":
			tokenResp := &TokenResponse{
				AccessToken: values.Get("access_token"),
				Scope:       values.Get("scope"),
				TokenType:   values.Get("token_type"),
			}
			return tokenResp, nil
		case values.Get("error") == "authorization_pending":
			continue POLL_LOOP
		default:
			return nil, fmt.Errorf("error: %s, description: %s, uri: %s",
				values.Get("error"), values.Get("error_description"), values.Get("error_uri"))
		}
	}
}

func init() {
	loginCommand.Flags().StringVarP(&loginFlags.Host, "host", "", "", "Host name to login")
	authCommand.AddCommand(loginCommand)
}
