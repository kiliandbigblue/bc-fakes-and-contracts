package bigcommerce

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/bigbluedisco/bc-fakes-and-contracts/domain/apiclient"
	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
)

// CreateWebhook creates a new webhook subscription in the BigCommerce API.
func (c *Client) CreateWebhook(ctx context.Context, req *apiclient.UpsertWebhookRequest) (*apiclient.UpsertWebhookResponse, error) {
	// Builds the url to call.
	createWebhookEndpoint := fmt.Sprintf("%s/stores/%s/v3/hooks", c.BaseURL(), c.storeHash)

	uwr := &apiclient.UpsertWebhookResponse{}
	res, err := c.request(ctx, http.MethodPost, createWebhookEndpoint, true, req, &uwr)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.StatusCode > 299 {
		return nil, &UnexpectedStatusError{StatusCode: res.StatusCode}
	}

	return uwr, nil
}

// Update an existing webhook subscription.
func (c *Client) UpdateWebhook(ctx context.Context, webhookID int, req *apiclient.UpsertWebhookRequest) (*apiclient.UpsertWebhookResponse, error) {
	// Builds the url to call.
	updateWebhookEndpoint := fmt.Sprintf("%s/stores/%s/v3/hooks/%d", c.BaseURL(), c.storeHash, webhookID)

	uwr := &apiclient.UpsertWebhookResponse{}
	res, err := c.request(ctx, http.MethodPut, updateWebhookEndpoint, true, req, &uwr)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.StatusCode > 299 {
		return nil, &UnexpectedStatusError{StatusCode: res.StatusCode}
	}

	return uwr, nil
}

// ListWebhooks calls Bigcommerce API to retrieve a paginated list of webhooks.
// After an initial filtering request has been executed, use the PageLink field to fetch the next page
// of results by setting it to the value retrieved in the Response.Meta.Pagination.Links.Next field.
func (c *Client) ListWebhooks(ctx context.Context, opts *apiclient.ListWebhooksOptions) (*apiclient.ListWebhooksResponse, error) {
	// Builds the url to call.
	listWebhooksEndpoint := fmt.Sprintf("%s/stores/%s/v3/hooks", c.BaseURL(), c.storeHash)

	switch {
	case opts != nil && opts.PageLink != "":
		// If PageLink is specified, use it
		listWebhooksEndpoint += opts.PageLink
	case opts != nil:
		// Encode the opts in the query.
		v, err := query.Values(opts)
		if err != nil {
			return nil, errors.Wrapf(err, "opts=%v", opts)
		}
		listWebhooksEndpoint += "?" + v.Encode()
		listWebhooksEndpoint = strings.TrimSuffix(listWebhooksEndpoint, "?")
	}

	lwr := &apiclient.ListWebhooksResponse{}
	res, err := c.request(ctx, http.MethodGet, listWebhooksEndpoint, true, nil, &lwr)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.StatusCode > 299 {
		return nil, &UnexpectedStatusError{StatusCode: res.StatusCode}
	}

	return lwr, nil
}
