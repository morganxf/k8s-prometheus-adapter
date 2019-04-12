package restmapper

import (
	"sync"

	"k8s.io/client-go/discovery"

	"gitlab.alipay-inc.com/antcloud-aks/aks-k8s-api/pkg/multitenancy"
	"gitlab.alipay-inc.com/antcloud-aks/aks-k8s-api/pkg/multitenancy/meta"
)

func (d *DeferredDiscoveryRESTMapper) ShallowCopyWithTenant(tenant multitenancy.TenantInfo) interface{} {
	copied := *d
	copied.initMu = sync.Mutex{}
	copied.cl = d.cl.(meta.TenantWise).ShallowCopyWithTenant(tenant).(discovery.CachedDiscoveryInterface)
	return &copied
}
