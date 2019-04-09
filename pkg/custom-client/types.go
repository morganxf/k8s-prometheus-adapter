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

// Package client provides bindings to the Prometheus HTTP API:
// http://prometheus.io/docs/querying/api/
package client

import (
	"net/url"
	"time"

	"github.com/prometheus/common/model"
)

const (
	KeyReservedPrefix     = "__"
	KeyTenant             = KeyReservedPrefix + "tenant_name"
	KeyWorkspace          = KeyReservedPrefix + "workspace_name"
	KeyMetricResourceType = KeyReservedPrefix + "metric.resource_type"
)

var queryArgsMap = map[string]string{
	KeyTenant:    "tenantName",
	KeyWorkspace: "workspaceName",
}

// APIResponse represents the response return by the API(metrics server).
type APIResponse struct {
	Data APIResponseData `json:"datas"`
}

type APIResponseData struct {
	Data APIResponseData2 `json:"querydata2"`
}

type APIResponseData2 struct {
	Metrics     []Metric `json:"datas"`
	IsSuccessed bool     `json:"success"`
	Massage     string   `json:"errMsg"`
}

// APIQueryOptions represents the options used by the API.
type APIQueryOptions struct {
	MetricName  string
	Labels      map[string]string
	Start       int64
	End         int64
	Step        time.Duration
	Aggregator  string
	QueryValues url.Values
}

// AddQueryValues add the value to key.
func (opt *APIQueryOptions) AddQueryValues(key, value string) {
	if opt.QueryValues == nil {
		opt.QueryValues = url.Values{}
	}
	opt.QueryValues.Add(key, value)
}

// Range represents a sliced time range with increments.
type Range struct {
	// Start and End are the boundaries of the time range.
	Start, End model.Time
	// Step is the maximum time between two slices within the boundaries.
	Step time.Duration
}

// QueryResult is the result if a query.
type QueryResult struct {
	Metrics Metric
}

type Metric struct {
	Name       string             `json:"metric"`
	Labels     map[string]string  `json:"tags"`
	DataPoints map[string]float64 `json:"dps"`
}
