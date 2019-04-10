package multitenant

import (
	"context"
	"fmt"

	"k8s.io/apiserver/pkg/endpoints/request"
)

// In multi-tenant mode, tenant info are injected into user's "extra" (which is a map).
// The following are the keys defining how to get the tenant info from a user object.
const (
	UserExtraTenantName    = "antcloud-aks-tenant-id"
	UserExtraWorkspaceName = "antcloud-aks-workspace-id"
	UserExtraClusterName   = "antcloud-aks-cluster-id"
)

// key used by tenant info map.
const (
	MapKeyTenantName    = "tenant_name"
	MapKeyWorkspaceName = "workspace_name"
	MapKeyClusterName   = "cluster_name"
)

// TenantInfo contains multi-tenant info.
type TenantInfo struct {
	TenantName    string
	WorkspaceName string
	ClusterName   string
}

// IsMultiTenant return true if TenantName and WorkspaceName are not empty.
func (t *TenantInfo) IsMultiTenant() bool {
	if t.TenantName != "" && t.WorkspaceName != "" {
		return true
	}
	return false
}

// ToMap converts TenantInfo to TenantInfoMap.
func (t *TenantInfo) ToMap() map[string]string {
	tenantInfoMap := map[string]string{
		MapKeyTenantName:    t.TenantName,
		MapKeyWorkspaceName: t.WorkspaceName,
		MapKeyClusterName:   t.ClusterName,
	}
	return tenantInfoMap
}

// GetTenantInfoFromContext extracts tenant info from the context.
func GetTenantInfoFromContext(ctx context.Context) (TenantInfo, error) {
	tenantInfo := TenantInfo{}
	userInfo, ok := request.UserFrom(ctx)
	if !ok {
		return tenantInfo, fmt.Errorf("unable to get user info from the context")
	}
	extra := userInfo.GetExtra()
	if extra == nil {
		return tenantInfo, fmt.Errorf("failed to extract tenant info from userInfo: nil extra")
	}
	if value, ok := extra[UserExtraTenantName]; ok && len(value) == 1 {
		tenantInfo.TenantName = value[0]
	}
	if value, ok := extra[UserExtraWorkspaceName]; ok && len(value) == 1 {
		tenantInfo.WorkspaceName = value[0]
	}
	if value, ok := extra[UserExtraClusterName]; ok && len(value) == 1 {
		tenantInfo.ClusterName = value[0]
	}
	if !tenantInfo.IsMultiTenant() {
		return tenantInfo, fmt.Errorf("insufficient multi-tenant info, tenantInfo: %+v", tenantInfo)
	}
	return tenantInfo, nil
}
