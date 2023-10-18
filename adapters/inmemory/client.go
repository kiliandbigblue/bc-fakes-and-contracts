package inmemory

import (
	"context"
	"errors"

	"github.com/bigbluedisco/bc-fakes-and-contracts/domain/apiclient"
)

type Client struct {
	i        int
	webhooks map[ /* scope*/ string]apiclient.Webhook
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) CreateWebhook(ctx context.Context, req *apiclient.UpsertWebhookRequest) (*apiclient.UpsertWebhookResponse, error) {
	return nil, errors.New("not implemented")
}

func (c *Client) UpdateWebhook(ctx context.Context, webhookID int, req *apiclient.UpsertWebhookRequest) (*apiclient.UpsertWebhookResponse, error) {
	return nil, errors.New("not implemented")
}

func (c *Client) ListWebhooks(ctx context.Context, opts *apiclient.ListWebhooksOptions) (*apiclient.ListWebhooksResponse, error) {
	return nil, errors.New("not implemented")
}

func (c *Client) DeleteWebhook(ctx context.Context, webhookID int) (*apiclient.DeleteWebhookResponse, error) {
	return nil, errors.New("not implemented")
}
