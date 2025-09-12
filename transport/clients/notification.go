package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type CreateNotification struct {
	UserID        string            `json:"user_id"`
	Type          string            `json:"type"`
	ReferenceType string            `json:"reference_type"`
	ReferenceID   string            `json:"reference_id"`
	Title         string            `json:"title"`
	Message       string            `json:"message"`
	Headers       map[string]string `json:"headers,omitempty"`
	Priority      string            `json:"priority"`
}

type NotifClient struct {
	base string
	http *http.Client
}

func NewNotifClient(baseURL string) *NotifClient {
	return &NotifClient{
		base: baseURL,
		http: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *NotifClient) Send(ctx context.Context, n CreateNotification) error {
	b, _ := json.Marshal(n)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.base+"/api/v1/notifications", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	uid := n.Headers["x-user-id"]
	if uid == "" {
		return fmt.Errorf("missing required header x-user-id")
	}

	for k, v := range n.Headers {
		req.Header.Set(k, v)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("notif service status %d", resp.StatusCode)
	}
	return nil
}
