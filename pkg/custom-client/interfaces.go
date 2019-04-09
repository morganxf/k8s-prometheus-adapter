/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package client

import (
	"context"
)

// Client is a MonitorServer client for the MonitorServer HTTP API.
// The "timeout" parameter for the HTTP API is set based on the context's deadline,
// when present and applicable.
type Client interface {
	// QueryRange runs a range query at the given time.
	QueryRange(ctx context.Context, queryOpts APIQueryOptions) (QueryResult, error)
}
