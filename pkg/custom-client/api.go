// Copyright 2017 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package client provides bindings to the MonitorServer HTTP API:
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

const queryRangeURL = "/private_api/metric/dataQuery"

// NewClient creates a Client for the given HTTP client and base URL (the location of the Prometheus server).
func NewClient(client *http.Client, baseURL *url.URL) Client {
	genericClient := NewGenericAPIClient(client, baseURL)
	return NewClientForAPI(genericClient)
}

// NewClientForAPI creates a Client for the given generic Prometheus API client.
func NewClientForAPI(client GenericAPIClient) Client {
	return &queryClient{
		api: client,
	}
}

// GenericAPIClient is a raw client to do http request.
// It knows how to appropriately deal with generic MonitorServer API
// responses, but does not know the specifics of different endpoints.
// You can use this to call query endpoints not represented in Client.
type GenericAPIClient interface {
	// Do makes a request to the Monitor-Server HTTP API against a particular endpoint.  Query
	// parameters should be in `query`, not `endpoint`.  An error will be returned on HTTP
	// status errors or errors making or unmarshalling the request, as well as when the
	// response has a Status of ResponseError.
	Do(ctx context.Context, verb, urlPath string, query url.Values, payload io.Reader) (APIResponse, error)
}

// NewGenericAPIClient builds a new generic Prometheus API client for the given base URL and HTTP Client.
func NewGenericAPIClient(client *http.Client, baseURL *url.URL) GenericAPIClient {
	return &httpAPIClient{
		client:  client,
		baseURL: baseURL,
	}
}

// httpAPIClient is a GenericAPIClient implemented in terms of an underlying http.Client.
type httpAPIClient struct {
	client  *http.Client
	baseURL *url.URL
}

func (c *httpAPIClient) Do(ctx context.Context, verb, urlPath string, query url.Values, payload io.Reader) (APIResponse, error) {
	var res APIResponse

	u := *c.baseURL
	u.Path = path.Join(c.baseURL.Path, urlPath)
	u.RawQuery = query.Encode()
	req, err := http.NewRequest(verb, u.String(), payload)
	if err != nil {
		return res, fmt.Errorf("error constructing HTTP request to MonitorServer: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.WithContext(ctx)

	resp, err := c.client.Do(req)
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	if err != nil {
		return res, fmt.Errorf("send request failed - %q, Error: %v", u.String(), err)
	}
	if resp.StatusCode/100 != 2 {
		return res, fmt.Errorf("request failed - %q", resp.Status)
	}

	// var body io.Reader = resp.Body
	// b, err := ioutil.ReadAll(body)
	// if err != nil {
	// 	return nil, fmt.Errorf("unable to read response body: %v", err)
	// }
	// body = bytes.NewReader(b)

	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return res, fmt.Errorf("failed to decode the output. Response: %q, Error: %v", err)
	}
	return res, nil
}

// queryClient is a Client that connects to the MonitorServer HTTP API.
type queryClient struct {
	api GenericAPIClient
}

// QueryRange implements Client interface.
func (c *queryClient) QueryRange(ctx context.Context, queryOpts APIQueryOptions) (QueryResult, error) {
	//return QueryResult{
	//	Metrics: Metric{
	//		DataPoints: map[string]float64{"1555060992": float64(60)},
	//	},
	//}, nil
	var queryRes QueryResult
	for _, key := range []string{KeyTenant, KeyWorkspace} {
		queryOpts.AddQueryValues(queryArgsMap[key], queryOpts.Labels[key])
	}

	queryBody := NewMonitorQueryBody(&queryOpts)
	b, err := json.Marshal(queryBody)
	if err != nil {
		return queryRes, err
	}
	res, err := c.api.Do(ctx, "POST", queryRangeURL, queryOpts.QueryValues, bytes.NewBuffer(b))
	if err != nil {
		return queryRes, err
	}

	if len(res.Data.Data.Metrics) == 0 {
		return queryRes, fmt.Errorf("not found. MetricName: %s, MetricLables: %+v", queryOpts.MetricName, queryOpts.Labels)
	}
	queryRes.Metrics = res.Data.Data.Metrics[0]
	return queryRes, nil
}

// timeoutFromContext checks the context for a deadline and calculates a "timeout" duration from it,
// when present
func timeoutFromContext(ctx context.Context) (time.Duration, bool) {
	if deadline, hasDeadline := ctx.Deadline(); hasDeadline {
		return time.Now().Sub(deadline), true
	}

	return time.Duration(0), false
}

// MonitorQueryBody is the body used to request Monitor API.
type MonitorQueryBody struct {
	Query MonitorQueryBody2 `json:"querys"`
}

type MonitorQueryBody2 struct {
	QueryData MonitorQueryData `json:"querydata2"`
}

type MonitorQueryData struct {
	MetricName   string            `json:"metricName"`
	Labels       map[string]string `json:"whiteTags"`
	Aggregator   string            `json:"aggregator"`
	Start        int64             `json:"startTime"`
	End          int64             `json:"endTime"`
	AttachAttr   map[string]string `json:"attachAttr"`
	ResourceType string            `json:"resourceType"`
}

func (m *MonitorQueryData) handleLabels() {
	for key := range m.Labels {
		if strings.HasPrefix(key, KeyReservedPrefix) {
			delete(m.Labels, key)
		}
	}
}

// NewMonitorQueryBody creates a MonitorQueryBody for the given APIQueryOptions.
func NewMonitorQueryBody(queryOpts *APIQueryOptions) MonitorQueryBody {
	const (
		defaultInterval   = -10 * time.Minute
		defaultAggregator = "none"
	)
	queryData := MonitorQueryData{
		MetricName:   queryOpts.MetricName,
		Labels:       DeepCopyLabels(queryOpts.Labels),
		Aggregator:   queryOpts.Aggregator,
		Start:        queryOpts.Start,
		End:          queryOpts.End,
		ResourceType: queryOpts.Labels[KeyMetricResourceType],
		AttachAttr:   make(map[string]string),
	}
	queryData.handleLabels()
	now := time.Now()
	if queryData.End == 0 {
		queryData.End = now.Unix() * 1000
	}
	if queryData.Start == 0 {
		queryData.Start = now.Add(defaultInterval).Unix() * 1000
	}
	if queryData.Aggregator == "" {
		queryData.Aggregator = defaultAggregator
	}

	queryBody := MonitorQueryBody{
		Query: MonitorQueryBody2{
			QueryData: queryData,
		},
	}
	return queryBody
}

// DeepCopyLabels returns a new labels.
func DeepCopyLabels(labels map[string]string) map[string]string {
	newLabels := make(map[string]string, len(labels))
	for key, value := range labels {
		newLabels[key] = value
	}
	return newLabels
}
