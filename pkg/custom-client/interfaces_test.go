package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	baseURLStr := "http://localhost:8080"
	baseURL, _ := url.Parse(baseURLStr)
	c := NewClient(http.DefaultClient, baseURL)
	queryOpts := APIQueryOptions{
		MetricName: "my-metric",
		Labels: map[string]string{
			"key-1": "value-1",
			// KeyTenant:    "tenant_name",
			// KeyWorkspace: "workspace_name",
		},
		Start:      time.Now().Unix(),
		End:        time.Now().Unix(),
		Step:       1 * time.Minute,
		Aggregator: "avg",
	}
	r, err := c.QueryRange(context.Background(), queryOpts)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("result: %+v\n", *r.Metrics)
}
