package apiclient

import (
	"context"
	"testing"

	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"
)

// Webhook contains the information BigCommerce provides about a webhook object.
type Webhook struct {
	// ID of the webhook.
	ID int `json:"id"`
	// Client ID, unique to the store or app.
	ClientID string `json:"client_id"`
	// Permanent ID of the BigCommerce store.
	StoreHash string `json:"store_hash"`
	// Event subscribed to. Example: "store/order/*".
	Scope string `json:"scope"`
	// Webhook destination URL. Must be active, return a 200 response,
	// and be served on port 443. Example: https://665b65a6.ngrok.io/webhooks.
	Destination string `json:"destination"`
	// Boolean value that indicates whether the webhook is active or not.
	IsActive bool `json:"is_active"`
	// Time when the webhook was created. Unix time.
	CreatedAt int `json:"created_at"`
	// Time when the webhook was last updated. Unix time.
	UpdatedAt int `json:"updated_at"`
	// Custom header.
	Headers map[string]string `json:"headers"`
}

// UpsertWebhookRequest represent the request body send to BigCommerce API "Create a Webhook" or "Update a Webhook".
type UpsertWebhookRequest struct {
	// Scope defines the event to subscribe to.
	Scope string `json:"scope"`
	// Destination is the URL where the webhook events are sent.
	Destination string `json:"destination"`
	// Active indicates whether the webhook is active or not.
	Active bool `json:"is_active"`
	// Custom header.
	Headers map[string]string `json:"headers"`
}

// UpsertWebhookResponse is the response retrieved from the BigCommerce API "Create a Webhook" or "Update a Webhook".
type UpsertWebhookResponse struct {
	// Webhook is the created webhook object.
	Webhook Webhook `json:"data"`
}

// ListWebhooksOptions is a struct used for querying a paginated collection of webhooks.
// It allows for initial filtering options to be set and supports paginated navigation.
type ListWebhooksOptions struct {
	// Page specifies the page number.
	Page int `url:"page,omitempty"`
	// Limit limits the number of results per page.
	Limit int `url:"limit,omitempty"`
	// Active is a filter for webhooks that are active or not.
	Active *bool `url:"is_active,omitempty"`
	// Scope is a filter for webhooks by scope.
	Scope string `url:"scope,omitempty"`
	// Destination is a filter for webhooks by destination.
	Destination string `url:"destination,omitempty"`

	// PageLink is used to store the link for fetching the next page of results.
	// When set, PageLink takes precedence over all other fields, which are effectively ignored.
	PageLink string `url:"-"`
}

// ListWebhooksResponse is the response retrieved from the Bigcommerce API Get Webhooks call.
type ListWebhooksResponse struct {
	// Data is the list of webhooks.
	Data []*Webhook `json:"data"`
	// Meta contains pagination links for the current, previous, and next parts of the whole collection.
	Meta *Meta `json:"meta,omitempty"`
}

// DeleteWebhookResponse is the response retrieved from the BigCommerce API Delete Webhook call.
type DeleteWebhookResponse struct {
	// Webhook is the deleted webhook object.
	Webhook Webhook `json:"data"`
}

type Client interface {
	CreateWebhook(ctx context.Context, req *UpsertWebhookRequest) (*UpsertWebhookResponse, error)
	UpdateWebhook(ctx context.Context, webhookID int, req *UpsertWebhookRequest) (*UpsertWebhookResponse, error)
	ListWebhooks(ctx context.Context, opts *ListWebhooksOptions) (*ListWebhooksResponse, error)
	DeleteWebhook(ctx context.Context, webhookID int) (*DeleteWebhookResponse, error)
}

type ClientContract struct {
	NewClient func() Client
}

func (c ClientContract) Test(t *testing.T) {
	t.Run("can create, get and list webooks", func(t *testing.T) {
		ctx := context.Background()
		cli := c.NewClient()

		scope := "store/order/updated"
		destination := "https://" + fake.DomainName() + "/" + fake.Characters()

		res, err := cli.CreateWebhook(ctx, &UpsertWebhookRequest{
			Scope:       scope,
			Destination: destination,
			Active:      true,
		})
		assert.NoError(t, err)

		got, err := cli.ListWebhooks(ctx, &ListWebhooksOptions{
			Scope: scope,
		})
		assert.NoError(t, err)
		assert.Len(t, got.Data, 1)
		assert.Equal(t, res, got.Data[0])

		newDestination := "https://" + fake.DomainName() + "/" + fake.Characters()
		res, err = cli.UpdateWebhook(ctx, res.Webhook.ID, &UpsertWebhookRequest{
			Destination: newDestination,
		})
		assert.NoError(t, err)

		got, err = cli.ListWebhooks(ctx, &ListWebhooksOptions{
			Scope: scope,
		})
		assert.NoError(t, err)
		assert.Len(t, got.Data, 1)
		assert.Equal(t, res, got.Data[0])
	})

	t.Run("the system does not allow to create a second webhook with the same scope", func(t *testing.T) {
		ctx := context.Background()
		cli := c.NewClient()

		scope := "store/order/updated"
		destination := "https://" + fake.DomainName() + "/" + fake.Characters()

		_, err := cli.CreateWebhook(ctx, &UpsertWebhookRequest{
			Scope:       scope,
			Destination: destination,
			Active:      true,
		})
		assert.NoError(t, err)

		newDestination := "https://" + fake.DomainName() + "/" + fake.Characters()
		_, err = cli.CreateWebhook(ctx, &UpsertWebhookRequest{
			Scope:       scope,
			Destination: newDestination,
			Active:      true,
		})
		assert.Error(t, err)
	})
}
