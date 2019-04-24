package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	mclient "github.com/directxman12/k8s-prometheus-adapter/pkg/custom-client"
	resprov "github.com/directxman12/k8s-prometheus-adapter/pkg/resource-provider"
	"github.com/golang/glog"
	basecmd "github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/cmd"
	resmetrics "github.com/kubernetes-incubator/metrics-server/pkg/apiserver/generic"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apiserver/pkg/util/feature"
	"k8s.io/apiserver/pkg/util/logs"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"gitlab.alipay-inc.com/antcloud-aks/aks-k8s-api/pkg/multitenancy"
)

type Adapter struct {
	basecmd.AdapterBase

	// MonitorServerURL is the endpoint describing how to connect to MonitorServer.
	MonitorServerURL string

	// AdapterConfigFile points to the file containing the configuration.
	AdapterConfigFile string
	KubeConfigFile    string

	// MetricsRelistInterval is the interval at which to relist the set of available metrics
	MetricsRelistInterval time.Duration
	// MetricsMaxAge is the period to query available metrics for
	MetricsMaxAge time.Duration
}

func (cmd *Adapter) addFlags() {
	cmd.Flags().StringVar(&cmd.MonitorServerURL, "monitor-server-url", "",
		"URL for connecting to MonitorServer.")
	cmd.Flags().StringVar(&cmd.AdapterConfigFile, "config", "",
		"Configuration file for metrics APIServer.")
	cmd.Flags().StringVar(&cmd.KubeConfigFile, "kube-config", "",
		"The path to the kubeconfig used to connect to the Kubernetes API server and the Kubelets (defaults to in-cluster config)")
	cmd.Flags().DurationVar(&cmd.MetricsRelistInterval, "metrics-relist-interval", 10*time.Minute,
		"interval at which to re-list the set of all available metrics from MonitorServer")
	cmd.Flags().DurationVar(&cmd.MetricsMaxAge, "metrics-max-age", 20*time.Minute,
		"period for which to query the set of available metrics from MonitorServer")
}

func (cmd *Adapter) makeMonitorClient() (mclient.Client, error) {
	if cmd.MonitorServerURL == "" {
		cmd.MonitorServerURL = os.Getenv("MONITOR_SERVER_URL")
	}
	if cmd.MonitorServerURL == "" {
		return nil, fmt.Errorf("invalid MonitorServer URL: empty url")
	}
	baseURL, err := url.Parse(cmd.MonitorServerURL)
	if err != nil {
		return nil, fmt.Errorf("invalid MonitorServer URL %q: %v", baseURL, err)
	}

	var httpClient *http.Client
	httpClient = http.DefaultClient

	return mclient.NewClient(httpClient, baseURL), nil
}

func (cmd *Adapter) addResourceMetricsAPI(mClient mclient.Client, kubeClient kubernetes.Interface) error {
	mapper, err := cmd.RESTMapper()
	if err != nil {
		return fmt.Errorf("cmd.RESTMapper - %v", err)
	}

	provider, err := resprov.NewProvider(mClient, kubeClient, mapper)
	if err != nil {
		return fmt.Errorf("unable to construct resource metrics API provider: %v", err)
	}

	provCfg := &resmetrics.ProviderConfig{
		Node: provider,
		Pod:  provider,
	}
	informers, err := cmd.Informers()
	if err != nil {
		return fmt.Errorf("cmd.Informers - %v", err)
	}

	server, err := cmd.Server()
	if err != nil {
		return fmt.Errorf("cmd.Server - %v", err)
	}

	if err := resmetrics.InstallStorage(provCfg, informers.Core().V1(), server.GenericAPIServer); err != nil {
		return fmt.Errorf("resmetrics.InstallStorage - %v", err)
	}

	return nil
}

func (cmd *Adapter) makeKubeClient() (kubernetes.Interface, error) {
	var err error
	// set up the client config
	var kubeConfig *rest.Config

	if len(cmd.KubeConfigFile) == 0 {
		cmd.KubeConfigFile = os.Getenv("KUBE_CONFIG_FILE")
	}
	if len(cmd.KubeConfigFile) > 0 {
		loader := &clientcmd.ClientConfigLoadingRules{ExplicitPath: cmd.KubeConfigFile}
		clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loader, &clientcmd.ConfigOverrides{})
		kubeConfig, err = clientConfig.ClientConfig()
	} else {
		kubeConfig, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, fmt.Errorf("unable to construct lister client config: %v", err)
	}

	// set up kubernetes client
	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("unable tp construct lister client: %v", err)
	}
	return kubeClient, nil
}

func init() {
	err := feature.DefaultFeatureGate.Add(map[feature.Feature]feature.FeatureSpec{
		multitenancy.FeatureName: {
			Default:    true,
			PreRelease: feature.Alpha,
		},
	})
	if err != nil {
		panic(fmt.Sprintf("failed to set DefaultFeatureGate: %v", err))
	}
}

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	// set up flags
	cmd := &Adapter{}
	cmd.Name = "metrics-apiserver"
	cmd.addFlags()
	cmd.Flags().AddGoFlagSet(flag.CommandLine) // make sure we get the glog flags
	if err := cmd.Flags().Parse(os.Args); err != nil {
		glog.Fatalf("unable to parse flags: %v", err)
	}

	// make the MonitorServer client
	mClient, err := cmd.makeMonitorClient()
	if err != nil {
		glog.Fatalf("unable to construct MonitorServer client: %v", err)
	}

	// load the config
	// TODO

	// make the Kubernetes client
	kubeClient, err := cmd.makeKubeClient()
	if err != nil {
		glog.Fatalf("failed to make kube client: %v", err)
	}

	// attach resource metrics supprot
	if err := cmd.addResourceMetricsAPI(mClient, kubeClient); err != nil {
		glog.Fatalf("unable to install resource metrics API: %v", err)
	}

	http.HandleFunc("/healthz", healthHandler)
	go func() {
		glog.Fatal(http.ListenAndServe(":9004", nil))
	}()
	// run the server
	if err := cmd.Run(wait.NeverStop); err != nil {
		glog.Fatalf("unable to run metrics APIServer: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok\n")
}
