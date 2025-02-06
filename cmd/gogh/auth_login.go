package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/kyoh86/gogh/v2"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var loginFlags struct {
	Host     string
	User     string
	Password string
}

const clientID = "Ov23li6aEWIxek6F8P5L" // ここに正しいClient IDを設定

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

var loginCommand = &cobra.Command{
	Use:   "login",
	Short: "Login for the host and owner",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		oauthConfig := &oauth2.Config{
			ClientID: clientID,
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://github.com/login/device/code",
				TokenURL: "https://github.com/login/oauth/access_token",
			},
			Scopes: []string{"repo"},
		}

		// Request device code
		deviceCodeResp, err := requestDeviceCode(oauthConfig.ClientID)
		if err != nil {
			return fmt.Errorf("failed to request device code: %w", err)
		}

		fmt.Printf("Visit %s and enter the code: %s\n", deviceCodeResp.VerificationURI, deviceCodeResp.UserCode)

		// Poll for token
		tokenResp, err := pollForToken(oauthConfig, deviceCodeResp)
		if err != nil {
			return fmt.Errorf("failed to poll for token: %w", err)
		}

		tokens.Set(loginFlags.Host, loginFlags.User, tokenResp.AccessToken)

		fmt.Println("Login successful!")
		return nil
	},
}

func requestDeviceCode(clientID string) (*DeviceCodeResponse, error) {
	resp, err := http.PostForm("https://github.com/login/device/code", url.Values{
		"client_id": {clientID},
		"scope":     {"repo"},
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Device code response body: %s\n", string(body)) // デバッグ出力

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
	for {
		time.Sleep(time.Duration(deviceCodeResp.Interval*2) * time.Second)

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

		fmt.Printf("Token response body: %s\n", string(body)) // デバッグ出力

		values, err := url.ParseQuery(string(body))
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == http.StatusOK && values.Get("error") == "" {
			tokenResp := &TokenResponse{
				AccessToken: values.Get("access_token"),
				Scope:       values.Get("scope"),
				TokenType:   values.Get("token_type"),
			}
			return tokenResp, nil
		} else if values.Get("error") == "authorization_pending" {
			continue
		} else {
			return nil, fmt.Errorf("error: %s, description: %s, uri: %s",
				values.Get("error"), values.Get("error_description"), values.Get("error_uri"))
		}
	}
}

func init() {
	loginCommand.Flags().StringVarP(&loginFlags.Host, "host", "", gogh.DefaultHost, "Host name to login")
	loginCommand.Flags().StringVarP(&loginFlags.User, "user", "", "", "User name to login")
	loginCommand.Flags().StringVarP(&loginFlags.Password, "password", "", "", `Password or developer private token

You should generate personal access tokens with "Repository permissions":

- ✅ Read-only access to "Contents" and "Metadata"
- ✅ Read and write access to "Administration"`)
	authCommand.AddCommand(loginCommand)
}
