package inmemory

import (
	"context"
	"errors"

	"github.com/bigbluedisco/bc-fakes-and-contracts/domain/apiclient"
)

type Client struct {
	i        int
	webhooks map[int]*apiclient.Webhook
}

func NewClient() *Client {
	return &Client{
		webhooks: make(map[int]*apiclient.Webhook),
	}
}

func (c *Client) CreateWebhook(ctx context.Context, req *apiclient.UpsertWebhookRequest) (*apiclient.UpsertWebhookResponse, error) {
	active := true
	if req.Active != nil {
		active = *req.Active
	}

	for _, wh := range c.webhooks {
		if wh.Scope == req.Scope {
			return nil, errors.New("a webhook with this scope already exists")
		}
	}

	wh := &apiclient.Webhook{
		ID:          c.i,
		Scope:       req.Scope,
		Destination: req.Destination,
		IsActive:    active,
		Headers:     req.Headers,
	}
	c.webhooks[c.i] = wh
	c.i++
	return &apiclient.UpsertWebhookResponse{
		Webhook: *wh,
	}, nil
}

func (c *Client) UpdateWebhook(ctx context.Context, webhookID int, req *apiclient.UpsertWebhookRequest) (*apiclient.UpsertWebhookResponse, error) {
	wh := c.webhooks[webhookID]
	if req.Scope != "" {
		wh.Scope = req.Scope
	}
	if req.Destination != "" {
		wh.Destination = req.Destination
	}
	if req.Active != nil {
		wh.IsActive = *req.Active
	}
	if req.Headers != nil {
		wh.Headers = req.Headers
	}

	c.webhooks[webhookID] = wh

	return &apiclient.UpsertWebhookResponse{
		Webhook: *wh,
	}, nil
}

func (c *Client) ListWebhooks(ctx context.Context, opts *apiclient.ListWebhooksOptions) (*apiclient.ListWebhooksResponse, error) {
	res := &apiclient.ListWebhooksResponse{}
	for _, wh := range c.webhooks {
		if opts.Active != nil {
			if wh.IsActive != *opts.Active {
				continue
			}
		}
		if opts.Scope != "" {
			if wh.Scope != opts.Scope {
				continue
			}
		}
		if opts.Destination != "" {
			if wh.Destination != opts.Destination {
				continue
			}
		}
		res.Data = append(res.Data, wh)
	}
	return res, nil
}
