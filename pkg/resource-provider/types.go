package provider

import (
	client "github.com/directxman12/k8s-prometheus-adapter/pkg/custom-client"
	corev1 "k8s.io/api/core/v1"
)

var resourceNames = []corev1.ResourceName{corev1.ResourceCPU, corev1.ResourceMemory}

type DataPoint struct {
	Timestamp int64
	Value     float64
}

type LabelSet map[string]string

type Metric struct {
	LabelSet
	DataPoints []DataPoint
}

type resourceInfo struct {
	namespace string
	name      string
}

type resourceMetric struct {
	resourceInfo

	cpu []Metric
	mem []Metric
}

type clusterInfo struct {
	tenantName    string
	workspaceName string
	clusterName   string
}

type containerInfo struct {
	name string
	id   string
	resources corev1.ResourceRequirements
}

func (c containerInfo) GetCapacity() corev1.ResourceList {
	capacity := make(corev1.ResourceList, 2)
	if c.resources.Limits != nil {
		for key, value := range c.resources.Limits {
			capacity[key] = value
		}
	}
	if c.resources.Requests != nil {
		for _, key := range resourceNames {
			if _, found := capacity[key]; !found {
				capacity[key] = c.resources.Requests[key]
			}
		}
	}
	return capacity
}

type nodeInfo struct {
	ip string
	capacity corev1.ResourceList
}

func (n nodeInfo) GetCapacity() corev1.ResourceList {
	return n.capacity
}

type nodeResource struct {
	clusterInfo
	resourceInfo
	nodeIP string
}

func newNodeResource(name, ip, tenantName, workspaceName, clusterName string) *nodeResource {
	return &nodeResource{
		clusterInfo: clusterInfo{
			tenantName:    tenantName,
			workspaceName: workspaceName,
			clusterName:   clusterName,
		},
		resourceInfo: resourceInfo{
			name: name,
		},
		nodeIP: ip,
	}
}

type queryOptsBuilder interface {
	buildQueryOpts(metricName string) []*client.APIQueryOptions
}

func (r *nodeResource) buildQueryOpts(metricName string) []*client.APIQueryOptions {
	return []*client.APIQueryOptions{
		{
			MetricName: metricName,
			Labels: map[string]string{
				"mip":                        r.nodeIP,
				client.KeyTenant:             r.tenantName,
				client.KeyWorkspace:          r.workspaceName,
				client.KeyMetricResourceType: "MACHINE",
			},
		},
	}
}

type podResource struct {
	clusterInfo
	resourceInfo
	containers []containerInfo
}

func newPodResource(name, namespace string, containers []containerInfo, tenantName, workspaceName, clusterName string) *podResource {
	return &podResource{
		clusterInfo: clusterInfo{
			tenantName:    tenantName,
			workspaceName: workspaceName,
			clusterName:   clusterName,
		},
		resourceInfo: resourceInfo{
			namespace: namespace,
			name:      name,
		},
		containers: containers,
	}
}

func (r *podResource) buildQueryOpts(metricName string) []*client.APIQueryOptions {
	opts := make([]*client.APIQueryOptions, len(r.containers))
	for i, cont := range r.containers {
		opts[i] = &client.APIQueryOptions{
			MetricName: metricName,
			Labels: map[string]string{
				"id":                        cont.id,
				client.KeyTenant:             r.tenantName,
				client.KeyWorkspace:          r.workspaceName,
				client.KeyMetricResourceType: "CONTAINER_POD",
			},
		}
	}
	return opts
}
