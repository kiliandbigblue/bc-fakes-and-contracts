package bigcommerce

import (
	"net/http"
	"time"
)

const (
	DefaultBaseURL        = "https://api.bigcommerce.com"
	DefaultOAuthBaseUrl   = "https://login.bigcommerce.com"
	DefaultTimeout        = 10 * time.Second
	DefaultMaxErrorLength = 1000
)

type Client struct {
	hc             *http.Client
	baseURL        string
	oAuthBaseURL   string
	storeHash      string
	accessToken    string
	maxErrorLength int
}

type clientOptions struct {
	BaseURL        string
	OAuthBaseURL   string
	StoreHash      string
	AccessToken    string
	Timeout        time.Duration
	MaxErrorLength int
}

type ClientOption func(*clientOptions)

func WithStore(storeHash, accessToken string) ClientOption {
	return func(co *clientOptions) {
		co.StoreHash = storeHash
		co.AccessToken = accessToken
	}
}

func WithBaseURL(url string) ClientOption {
	return func(co *clientOptions) {
		co.BaseURL = url
	}
}

func WithOAuthBaseUrl(url string) ClientOption {
	return func(co *clientOptions) {
		co.OAuthBaseURL = url
	}
}

func WithTimeout(to time.Duration) ClientOption {
	return func(co *clientOptions) {
		co.Timeout = to
	}
}

func WithMaxErrorLength(limit int) ClientOption {
	return func(co *clientOptions) {
		co.MaxErrorLength = limit
	}
}

func NewClient(options ...ClientOption) *Client {
	opts := &clientOptions{
		BaseURL:        DefaultBaseURL,
		OAuthBaseURL:   DefaultOAuthBaseUrl,
		Timeout:        DefaultTimeout,
		MaxErrorLength: DefaultMaxErrorLength,
	}
	for _, o := range options {
		o(opts)
	}

	return &Client{
		hc: &http.Client{
			Timeout: opts.Timeout,
		},
		baseURL:        opts.BaseURL,
		oAuthBaseURL:   opts.OAuthBaseURL,
		storeHash:      opts.StoreHash,
		accessToken:    opts.AccessToken,
		maxErrorLength: opts.MaxErrorLength,
	}
}

func (c *Client) BaseURL() string {
	return c.baseURL
}

func (c *Client) OAuthBaseURL() string {
	return c.oAuthBaseURL
}
