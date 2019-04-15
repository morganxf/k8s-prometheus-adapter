package provider

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	client "github.com/directxman12/k8s-prometheus-adapter/pkg/custom-client"
	"github.com/directxman12/k8s-prometheus-adapter/pkg/multitenant"
	"github.com/golang/glog"
	"github.com/kubernetes-incubator/metrics-server/pkg/provider"
	apiv1 "k8s.io/api/core/v1"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/api/resource"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apitypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/metrics/pkg/apis/metrics"

	tenantmeta "gitlab.alipay-inc.com/antcloud-aks/aks-k8s-api/pkg/multitenancy/meta"
)

// NewProvider constructs a new MetricsProvider to provide resource metrics from MonitorServer.
func NewProvider(c client.Client, kubeClient kubernetes.Interface, mapper apimeta.RESTMapper) (provider.MetricsProvider, error) {
	cpuQuery := newResourceQuery("cpu_util")
	memQuery := newResourceQuery("mem_util")
	return &resourceProvider{
		client: c,
		cpu:    cpuQuery,
		mem:    memQuery,

		kubeClient: kubeClient,
	}, nil
}

// resourceProvider is a MetricsProvider that contacts MonitorServer to provide
// the resource metrics.
type resourceProvider struct {
	client client.Client
	cpu    resourceQuery
	mem    resourceQuery
	window time.Duration

	kubeClient kubernetes.Interface
}

// GetNodeMetrics implements the provider.MetricsProvider interface. It may return nil, nil, nil.
func (p *resourceProvider) GetNodeMetrics(ctx context.Context, nodes ...string) ([]provider.TimeInfo, []apiv1.ResourceList, error) {
	glog.Infof("start to get nodes metrics: %v", nodes)
	if len(nodes) == 0 {
		return nil, nil, nil
	}
	tenantInfo, err := multitenant.GetTenantInfoFromContext(ctx)
	if err != nil {
		glog.Errorf("failed to fetch tenant info for nodes %q...: %v", nodes[0], err)
	}

	kc, ok := p.kubeClient.(tenantmeta.TenantWise).ShallowCopyWithTenant(tenantInfo).(kubernetes.Interface)
	if !ok {
		glog.Error("type assertion failed")
		return nil, nil, nil
	}
	nodeIPs := make([]string, len(nodes))
	for i, nodeName := range nodes {
		// node, err := p.kubeClient.CoreV1().Nodes().Get(nodeName, meta_v1.GetOptions{})
		node, err := kc.CoreV1().Nodes().Get(nodeName, meta_v1.GetOptions{})
		if err != nil {
			glog.Errorf("failed to get node object %q: %v. continue", nodeName, err)
			continue
		}
		nodeIP, err := getNodeAddress(node)
		if err != nil {
			glog.Errorf("failed to get node address %q: %v. continue", nodeName, err)
			continue
		}
		nodeIPs[i] = nodeIP
	}

	builders := make([]queryOptsBuilder, len(nodes))
	for i, nodeName := range nodes {
		nr := newNodeResource(nodeName, nodeIPs[i], tenantInfo.TenantName, tenantInfo.WorkspaceName, tenantInfo.ClusterName)
		builders[i] = nr
	}

	rawResMetrics := p.queryBoth(builders...)

	resTimes := make([]provider.TimeInfo, len(nodes))
	resMetrics := make([]apiv1.ResourceList, len(nodes))

	// organize the results
	for i, nodeName := range nodes {
		// skip if any data is missing
		rm := rawResMetrics[i]
		if rm == nil || len(rm.cpu) == 0 || len(rm.mem) == 0 {
			glog.Infof("missing resource metrics for node %q, skipping", nodeName)
			continue
		}
		if len(rm.cpu[0].DataPoints) == 0 {
			glog.Infof("missing CPU metric for node %q, skipping", nodeName)
			continue
		}
		if len(rm.mem[0].DataPoints) == 0 {
			glog.Infof("missing memory metric for node %q, skipping", nodeName)
			continue
		}

		cpu := rm.cpu[0].DataPoints[0]
		mem := rm.mem[0].DataPoints[0]

		// store the results
		resMetrics[i] = apiv1.ResourceList{
			//apiv1.ResourceCPU:    *resource.NewMilliQuantity(int64(cpu.Value*1000.0), resource.DecimalSI),
			//apiv1.ResourceMemory: *resource.NewMilliQuantity(int64(mem.Value*1000.0), resource.BinarySI),
			apiv1.ResourceCPU:    *resource.NewMilliQuantity(int64(cpu.Value*1000.0), resource.DecimalExponent),
			apiv1.ResourceMemory: *resource.NewMilliQuantity(int64(mem.Value*1000.0), resource.DecimalExponent),
		}

		// use the earliest timestamp available (in order to be conservative
		// when determining if metrics are tainted by startup)
		if cpu.Timestamp < mem.Timestamp {
			resTimes[i] = provider.TimeInfo{
				Timestamp: time.Unix(cpu.Timestamp, 0),
				Window:    p.window,
			}
		} else {
			resTimes[i] = provider.TimeInfo{
				Timestamp: time.Unix(mem.Timestamp, 0),
				Window:    p.window,
			}
		}
	}

	glog.Infof("resTimes: %+v, resMetrics: %+v", resTimes, resMetrics)
	return resTimes, resMetrics, nil
}

// GetContainerMetrics implements the provider.MetricsProvider interface. It may return nil, nil, nil.
func (p *resourceProvider) GetContainerMetrics(ctx context.Context, pods ...apitypes.NamespacedName) ([]provider.TimeInfo, [][]metrics.ContainerMetrics, error) {
	glog.Infof("start to get pods metrics: %+v", pods)
	if len(pods) == 0 {
		return nil, nil, nil
	}
	tenantInfo, err := multitenant.GetTenantInfoFromContext(ctx)
	if err != nil {
		glog.Errorf("failed to fetch tenant info for pods \"%s/%s\"...: %v", pods[0].Namespace, pods[0].Name, err)
	}

	kc, ok := p.kubeClient.(tenantmeta.TenantWise).ShallowCopyWithTenant(tenantInfo).(kubernetes.Interface)
	if !ok {
		glog.Error("type assertion failed")
		return nil, nil, nil
	}

	podContainerInfos := make([][]containerInfo, len(pods))
	for i, podNameInfo := range pods {
		// pod, err := p.kubeClient.CoreV1().Pods(podNameInfo.Namespace).Get(podNameInfo.Name, meta_v1.GetOptions{})
		pod, err := kc.CoreV1().Pods(podNameInfo.Namespace).Get(podNameInfo.Name, meta_v1.GetOptions{})
		if err != nil {
			glog.Errorf("failed to get pod object \"%s/%s\": %v. continue", podNameInfo.Namespace, podNameInfo.Name, err)
			continue
		}
		containerInfos := make([]containerInfo, len(pod.Status.ContainerStatuses))
		for j, containerStatus := range pod.Status.ContainerStatuses {
			containerInfos[j] = containerInfo{
				name: containerStatus.Name,
				id:   strings.TrimLeft(containerStatus.ContainerID, "docker://"),
			}
		}
		podContainerInfos[i] = containerInfos
	}

	glog.Infof("podContainerInfos: %v", podContainerInfos)

	builders := make([]queryOptsBuilder, len(pods))
	for i, podNameInfo := range pods {
		podName := podNameInfo.Name
		podNamespace := podNameInfo.Namespace
		pr := newPodResource(podName, podNamespace, podContainerInfos[i], tenantInfo.TenantName, tenantInfo.WorkspaceName, tenantInfo.ClusterName)
		builders[i] = pr
	}

	rawResMetrics := p.queryBoth(builders...)

	resTimes := make([]provider.TimeInfo, len(pods))
	resMetrics := make([][]metrics.ContainerMetrics, len(pods))

	for i, podNameInfo := range pods {
		rm := rawResMetrics[i]
		if rm == nil {
			glog.Infof("missing resource metrics for pod \"%s/%s\", skipping", podNameInfo.Namespace, podNameInfo.Name)
			continue
		}

		earliestTs := time.Now().Unix()

		containerInfos := podContainerInfos[i]
		containerMetrics := make([]metrics.ContainerMetrics, len(containerInfos))
		for j := 0; j < len(containerInfos); j++ {
			containerInfo := containerInfos[j]
			if len(rm.cpu[j].DataPoints) == 0 {
				glog.Infof("missing CPU metric for pod container \"%s/%s/%s\", skipping", podNameInfo.Namespace, podNameInfo.Name, containerInfo.name)
				continue
			}
			if len(rm.mem[j].DataPoints) == 0 {
				glog.Infof("missing memroy metric for pod container \"%s/%s/%s\", skipping", podNameInfo.Namespace, podNameInfo.Name, containerInfo.name)
				continue
			}
			cpu := rm.cpu[j].DataPoints[0]
			mem := rm.mem[j].DataPoints[0]

			containerMetrics[j] = metrics.ContainerMetrics{
				Name: containerInfo.name,
				Usage: apiv1.ResourceList{
					//apiv1.ResourceCPU:    *resource.NewMilliQuantity(int64(cpu.Value*1000.0), resource.DecimalSI),
					//apiv1.ResourceMemory: *resource.NewMilliQuantity(int64(mem.Value*1000.0), resource.BinarySI),
					apiv1.ResourceCPU:    *resource.NewMilliQuantity(int64(cpu.Value*1000.0), resource.DecimalExponent),
					apiv1.ResourceMemory: *resource.NewMilliQuantity(int64(mem.Value*1000.0), resource.DecimalExponent),
				},
			}

			if cpu.Timestamp < earliestTs {
				earliestTs = cpu.Timestamp
			}
			if mem.Timestamp < earliestTs {
				earliestTs = mem.Timestamp
			}
		}

		resTimes[i] = provider.TimeInfo{
			Timestamp: time.Unix(earliestTs, 0),
			Window:    p.window,
		}

		resMetrics[i] = containerMetrics
	}

	glog.Infof("pods resTimes: %+v, resMetrics: %+v", resTimes, resMetrics)
	return resTimes, resMetrics, nil
}

func (p *resourceProvider) queryBoth(builders ...queryOptsBuilder) []*resourceMetric {
	var resCPUMetrics, resMemoryMetrics [][]Metric

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		resCPUMetrics = p.queryCPUMetrics(builders...)
	}()
	go func() {
		defer wg.Done()
		resMemoryMetrics = p.queryMemoryMetrics(builders...)
	}()
	wg.Wait()

	resMetrics := make([]*resourceMetric, len(builders))
	for i := 0; i < len(resMetrics); i++ {
		resMetrics[i] = &resourceMetric{
			cpu: resCPUMetrics[0],
			mem: resMemoryMetrics[0],
		}
	}

	return resMetrics
}

func (p *resourceProvider) queryCPUMetrics(builders ...queryOptsBuilder) [][]Metric {
	resQueryOpts := make([][]*client.APIQueryOptions, len(builders))
	for i, builder := range builders {
		resQueryOpts[i] = builder.buildQueryOpts("cpu_util")
	}

	resMetrics := p.queryResourceMetrics(resQueryOpts...)

	return resMetrics
}

func (p *resourceProvider) queryMemoryMetrics(builders ...queryOptsBuilder) [][]Metric {
	resQueryOpts := make([][]*client.APIQueryOptions, len(builders))
	for i, builder := range builders {
		resQueryOpts[i] = builder.buildQueryOpts("mem_util")
	}

	resMetrics := p.queryResourceMetrics(resQueryOpts...)

	return resMetrics
}

func (p *resourceProvider) queryResourceMetrics(resQueryOpts ...[]*client.APIQueryOptions) [][]Metric {
	resMetrics := make([][]Metric, 0, len(resQueryOpts))

	msChan := make(chan []Metric, len(resQueryOpts))
	var wg sync.WaitGroup
	wg.Add(len(resQueryOpts))
	for _, queryOpts := range resQueryOpts {
		go func(queryOpts []*client.APIQueryOptions) {
			defer wg.Done()
			result := p.queryMetrics(queryOpts)
			msChan <- result
		}(queryOpts)
	}
	wg.Wait()
	close(msChan)

	for result := range msChan {
		resMetrics = append(resMetrics, result)
	}

	return resMetrics
}

func (p *resourceProvider) queryMetrics(queryOpts []*client.APIQueryOptions) []Metric {
	memMetrics := make([]Metric, 0, len(queryOpts))

	resChan := make(chan *client.QueryResult, len(queryOpts))
	var wg sync.WaitGroup
	wg.Add(len(queryOpts))

	for _, queryOpt := range queryOpts {
		go func(ctx context.Context, queryOpt *client.APIQueryOptions) {
			defer wg.Done()
			result, err := p.client.QueryRange(ctx, *queryOpt)
			if err != nil {
				glog.Errorf("failed to query MonitorServer. queryOpts: %+v:, err: %v", queryOpts, err)
				resChan <- nil
				return
			}
			resChan <- &result
		}(context.TODO(), queryOpt)
	}

	wg.Wait()
	close(resChan)

	for result := range resChan {
		if result == nil {
			memMetrics = append(memMetrics, Metric{})
			continue
		}

		metric := Metric{
			DataPoints: make([]DataPoint, len(result.Metrics.DataPoints)),
		}

		// reverse sort timestamps
		timeStrs := make([]string, 0, len(result.Metrics.DataPoints))
		for k := range result.Metrics.DataPoints {
			timeStrs = append(timeStrs, k)
		}
		sort.Sort(sort.Reverse(sort.StringSlice(timeStrs)))

		for i, timeStr := range timeStrs {
			v := result.Metrics.DataPoints[timeStr]
			timestamp, err := strconv.ParseInt(timeStr, 10, 32)
			if err != nil {
				glog.Errorf("failed to prase timestamp %q: %v", timeStr, err)
				continue
			}
			metric.DataPoints[i] = DataPoint{
				Timestamp: timestamp,
				Value:     v,
			}
		}
		memMetrics = append(memMetrics, metric)
	}

	return memMetrics
}

// resourceQuery represents query information for querying resource metrics for some resource,
// like CPU or memory.
type resourceQuery struct {
	contQuery queryOpts
	nodeQuery queryOpts
}

func newResourceQuery(metricName string) resourceQuery {
	return resourceQuery{
		contQuery: queryOpts{
			metricName:   metricName,
			resourceType: "CONTAINER",
		},
		nodeQuery: queryOpts{
			metricName:   metricName,
			resourceType: "NODE",
		},
	}
}

type queryOpts struct {
	metricName   string
	resourceType string
}
