package provider

import (
	"fmt"

	apiv1 "k8s.io/api/core/v1"
)

const (
	NodeLegacyHostIP = "LegacyHostIP"
)

// getNodeAddress returns the provided node's address, based on the priority:
// 1. NodeInternalIP
// 2. NodeExternalIP
// 3. NodeLegacyHostIP
// 3. NodeHostName
//
// Derived from k8s.io/kubernetes/pkg/util/node/node.go
func getNodeAddress(node *apiv1.Node) (string, error) {
	m := map[apiv1.NodeAddressType][]string{}
	for _, a := range node.Status.Addresses {
		m[a.Type] = append(m[a.Type], a.Address)
	}

	if addresses, ok := m[apiv1.NodeInternalIP]; ok {
		return addresses[0], nil
	}
	if addresses, ok := m[apiv1.NodeExternalIP]; ok {
		return addresses[0], nil
	}
	if addresses, ok := m[apiv1.NodeAddressType(NodeLegacyHostIP)]; ok {
		return addresses[0], nil
	}
	if addresses, ok := m[apiv1.NodeHostName]; ok {
		return addresses[0], nil
	}
	return "", fmt.Errorf("failed to get node address")
}
