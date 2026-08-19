package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client is a thin wrapper over the Disassembly.AI REST API.
type Client struct {
	HTTP     *http.Client
	Endpoint string
	APIKey   string
}

type Finding struct {
	RuleID   string `json:"rule_id"`
	Title    string `json:"title"`
	Severity string `json:"severity"` // "note" | "warning" | "error"
}

type ScanResult struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"` // "queued" | "running" | "completed" | "failed"
	Findings  []Finding `json:"findings"`
	SARIF     string    `json:"sarif"`
	ReportURL string    `json:"report_url"`
	Tokens    int64     `json:"tokens"`
	Runs      int64     `json:"runs"`
}

func (c *Client) do(ctx context.Context, method, path string, body any, out any) error {
	var rdr *bytes.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		rdr = bytes.NewReader(b)
	} else {
		rdr = bytes.NewReader(nil)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.Endpoint+path, rdr)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	res, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		return fmt.Errorf("%s %s: %s", method, path, res.Status)
	}
	return json.NewDecoder(res.Body).Decode(out)
}

func (c *Client) StartScan(ctx context.Context, target, effort string) (*ScanResult, error) {
	var out ScanResult
	err := c.do(ctx, http.MethodPost, "/v1/scans",
		map[string]string{"target": target, "effort": effort}, &out)
	return &out, err
}

func (c *Client) GetScan(ctx context.Context, id string) (*ScanResult, error) {
	var out ScanResult
	err := c.do(ctx, http.MethodGet, "/v1/scans/"+id, nil, &out)
	return &out, err
}

// WaitForScan polls until the scan reaches a terminal status.
func (c *Client) WaitForScan(ctx context.Context, id string) (*ScanResult, error) {
	for {
		r, err := c.GetScan(ctx, id)
		if err != nil {
			return nil, err
		}
		switch r.Status {
		case "completed":
			return r, nil
		case "failed":
			return nil, fmt.Errorf("scan %s failed", id)
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(5 * time.Second):
		}
	}
}
