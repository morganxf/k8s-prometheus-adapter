package discovery

import (
	"sync"

	"gitlab.alipay-inc.com/antcloud-aks/aks-k8s-api/pkg/multitenancy"
	"gitlab.alipay-inc.com/antcloud-aks/aks-k8s-api/pkg/multitenancy/meta"
)

func (d *CachedDiscoveryClient) ShallowCopyWithTenant(tenant multitenancy.TenantInfo) interface{} {
	copied := *d
	copied.mutex = sync.Mutex{}
	copied.delegate = d.delegate.(meta.TenantWise).ShallowCopyWithTenant(tenant).(DiscoveryInterface)
	return &copied
}
