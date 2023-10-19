package bigcommerce

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// apiResponse represents a http response.
type apiResponse struct {
	Header http.Header

	// RawBody contains the response body as raw bytes.
	RawBody []byte

	// StatusCode is a status code as integer. e.g. 200.
	StatusCode int
}

// ResponseMetadata represent pagination and collection totals for multi-page responses.
type Meta struct {
	Pagination Pagination `json:"pagination"`
}

// Pagination and collection totals in the response.
type Pagination struct {
	Total       int   `json:"total"`
	Count       int   `json:"count"`
	PerPage     int   `json:"per_page"`
	CurrentPage int   `json:"current_page"`
	TotalPages  int   `json:"total_pages"`
	Links       Links `json:"links"`
}

// Pagination links for the current, previous and next parts of the whole collection.
type Links struct {
	Previous string `json:"previous"`
	Current  string `json:"current"`
	Next     string `json:"next"`
}

// Error returned when the api response status does not match the expected one.
type UnexpectedStatusError struct {
	StatusCode int
}

func (e *UnexpectedStatusError) Error() string {
	return fmt.Sprintf("wrong status - got: %d", e.StatusCode)
}

// ErrNoContentToUnmarshal is returned when the API response body is empty, but the request expected to unmarshal data from it.
var ErrNoContentToUnmarshal = errors.New("no content to unmarshal")

// request calls the url endpoint using the client authentication information
// It returns an APIResponse that contains the request response information
// It returns an error in case of marshaling or request issue.
func (c *Client) request(ctx context.Context, method, url string, authRequired bool, reqBody, resUnmarshaled interface{}) (*apiResponse, error) {
	var jsonReqBody []byte
	var err error

	if reqBody != nil {
		jsonReqBody, err = json.Marshal(reqBody)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(jsonReqBody))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	req.Header.Add("Accept", "application/json")

	if reqBody != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	if authRequired {
		req.Header.Add("X-Auth-Token", c.accessToken)
	}

	res, err := c.hc.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() { _ = res.Body.Close() }()

	rawResBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	switch {
	// Error responses.
	case res.StatusCode > 399:
		body := string(rawResBody)
		if len(body) > c.maxErrorLength {
			body = body[:c.maxErrorLength]
		}
		return nil, errors.Errorf("%s - body: %s", res.Status, body)

	case res.StatusCode < 300 && resUnmarshaled != nil && len(rawResBody) == 0:
		// The response is successful, but its body is empty,
		// which is unexpected since the request was supposed to unmarshal data from it.
		return &apiResponse{
			Header:     res.Header,
			StatusCode: res.StatusCode,
		}, ErrNoContentToUnmarshal

	case res.StatusCode < 300 && resUnmarshaled != nil:
		// Successful responses.
		err := json.Unmarshal(rawResBody, resUnmarshaled)
		if err != nil {
			body := string(rawResBody)
			if len(body) > c.maxErrorLength {
				body = body[:c.maxErrorLength]
			}
			return nil, errors.Wrapf(err, "%s - body: %s", res.Status, body)
		}
	}

	return &apiResponse{
		Header:     res.Header,
		RawBody:    rawResBody,
		StatusCode: res.StatusCode,
	}, nil
}
