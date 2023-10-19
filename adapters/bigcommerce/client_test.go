package bigcommerce

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	_fakeBaseURL = "https://fake.bigcommerce.com"
)

func TestNewClient_OK_Default(t *testing.T) {
	c := NewClient()

	assert.Equal(t, DefaultTimeout, c.hc.Timeout)
	assert.Equal(t, DefaultBaseURL, c.baseURL)
	assert.Equal(t, DefaultOAuthBaseUrl, c.oAuthBaseURL)
	assert.Equal(t, "", c.storeHash)
	assert.Equal(t, "", c.accessToken)
	assert.Equal(t, DefaultMaxErrorLength, c.maxErrorLength)
}

func TestNewClient_OK_Options(t *testing.T) {
	timeout := 15 * time.Second
	hash := "store_hash"
	token := "access_token"
	url := _fakeBaseURL
	limit := 100
	c := NewClient(WithTimeout(timeout), WithStore(hash, token), WithBaseURL(url), WithOAuthBaseUrl(url), WithMaxErrorLength(limit))

	assert.Equal(t, timeout, c.hc.Timeout)
	assert.Equal(t, url, c.baseURL)
	assert.Equal(t, url, c.oAuthBaseURL)
	assert.Equal(t, hash, c.storeHash)
	assert.Equal(t, token, c.accessToken)
	assert.Equal(t, limit, c.maxErrorLength)
}

func TestBaseURL_OK(t *testing.T) {
	url := _fakeBaseURL
	c := NewClient(WithBaseURL(url))
	assert.Equal(t, url, c.BaseURL())
}

func TestOAuthBaseURL_OK(t *testing.T) {
	url := _fakeBaseURL
	c := NewClient(WithOAuthBaseUrl(url))
	assert.Equal(t, url, c.OAuthBaseURL())
}
