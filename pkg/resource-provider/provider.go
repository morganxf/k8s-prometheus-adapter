package provider

import (
	"context"
	"time"

	client "github.com/directxman12/k8s-prometheus-adapter/pkg/custom-client"
	"github.com/directxman12/k8s-prometheus-adapter/pkg/multitenant"
	"github.com/golang/glog"
	"github.com/kubernetes-incubator/metrics-server/pkg/provider"
	corev1 "k8s.io/api/core/v1"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	apitypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/metrics/pkg/apis/metrics"
)

// NewProvider constructs a new MetricsProvider to provide resource metrics from MonitorServer.
func NewProvider(c client.Client, mapper apimeta.RESTMapper) (provider.MetricsProvider, error) {
	cpuQuery := newResourceQuery("cpu_util")
	memQuery := newResourceQuery("mem_util")
	return &resourceProvider{
		client: c,
		cpu:    cpuQuery,
		mem:    memQuery,
	}, nil
}

// resourceProvider is a MetricsProvider that contacts MonitorServer to provide
// the resource metrics.
type resourceProvider struct {
	client client.Client
	cpu    resourceQuery
	mem    resourceQuery
	window time.Duration
}

// GetNodeMetrics implements the provider.MetricsProvider interface. It may return nil, nil, nil.
func (p *resourceProvider) GetNodeMetrics(ctx context.Context, nodes ...string) ([]provider.TimeInfo, []corev1.ResourceList, error) {
	if len(nodes) == 0 {
		return nil, nil, nil
	}
	_, err := multitenant.GetTenantInfoFromContext(ctx)
	if err != nil {
		glog.Errorf("failed to fetch tenant info for nodes %q...: %v", nodes[0], err)
	} else {

	}
	return nil, nil, nil
}

// GetContainerMetrics implements the provider.MetricsProvider interface. It may return nil, nil, nil.
func (p *resourceProvider) GetContainerMetrics(ctx context.Context, pods ...apitypes.NamespacedName) ([]provider.TimeInfo, [][]metrics.ContainerMetrics, error) {
	if len(pods) == 0 {
		return nil, nil, nil
	}
	_, err := multitenant.GetTenantInfoFromContext(ctx)
	if err != nil {
		glog.Errorf("failed to fetch tenant info for pods \"%s/%s\"...: %v", pods[0].Namespace, pods[0].Name, err)
	} else {

	}
	return nil, nil, nil
}

// resourceQuery represents query information for querying resource metrics for some resource,
// like CPU or memory.
type resourceQuery struct {
	contQuery client.APIQueryOptions
	nodeQuery client.APIQueryOptions
}

func newResourceQuery(metricName string) resourceQuery {
	return resourceQuery{
		contQuery: client.APIQueryOptions{
			MetricName: metricName,
			Labels: map[string]string{
				client.KeyMetricResourceType: "CONTAINER",
			},
		},
		nodeQuery: client.APIQueryOptions{
			MetricName: metricName,
			Labels: map[string]string{
				client.KeyMetricResourceType: "NODE",
			},
		},
	}
}
