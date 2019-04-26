#!/bin/sh


curl -k --cert ~/kubernetes/aks/cert.pem --key ~/kubernetes/aks/key.pem \
-H "X-Remote-Extra-Antcloud-Aks-Cluster-Id: sgtest" \
-H "X-Remote-Extra-Antcloud-Aks-Tenant-Id: CLDUSGCN" \
-H "X-Remote-Extra-Antcloud-Aks-Workspace-Id: linkepre" \
-H "X-Remote-Group: system:masters" \
-H "X-Remote-User: system:apiserver" \
-H "User-Agent: hyperkube/v1.12.0 (darwin/amd64) kubernetes/31dda1c" \
-H "X-Forwarded-For: ::1" \
-H "Accept-Encoding: gzip" \
-H "Accept: application/vnd.kubernetes.protobuf, */*"  \
https://localhost:443/apis/metrics.k8s.io/v1beta1/nodes


curl -k --cert ~/kubernetes/aks/cert.pem --key ~/kubernetes/aks/key.pem \
-H "X-Remote-Extra-Antcloud-Aks-Cluster-Id: sgtest" \
-H "X-Remote-Extra-Antcloud-Aks-Tenant-Id: CLDUSGCN" \
-H "X-Remote-Extra-Antcloud-Aks-Workspace-Id: linkepre" \
-H "X-Remote-Group: system:masters" \
-H "X-Remote-User: system:apiserver" \
-H "User-Agent: hyperkube/v1.12.0 (darwin/amd64) kubernetes/31dda1c" \
-H "X-Forwarded-For: ::1" \
-H "Accept-Encoding: gzip" \
-H "Accept: application/vnd.kubernetes.protobuf, */*"  \
https://localhost:443/apis/metrics.k8s.io/v1beta1/nodes/2122587225


curl -k --cert ~/kubernetes/aks/cert.pem --key ~/kubernetes/aks/key.pem \
-H "X-Remote-Extra-Antcloud-Aks-Cluster-Id: sgtest" \
-H "X-Remote-Extra-Antcloud-Aks-Tenant-Id: CLDUSGCN" \
-H "X-Remote-Extra-Antcloud-Aks-Workspace-Id: linkepre" \
-H "X-Remote-Group: system:masters" \
-H "X-Remote-User: system:apiserver" \
-H "User-Agent: hyperkube/v1.12.0 (darwin/amd64) kubernetes/31dda1c" \
-H "X-Forwarded-For: ::1" \
-H "Accept-Encoding: gzip" \
-H "Accept: application/vnd.kubernetes.protobuf, */*"  \
https://localhost:443/apis/metrics.k8s.io/v1beta1/pods


curl -k --cert ~/kubernetes/aks/cert.pem --key ~/kubernetes/aks/key.pem \
-H "X-Remote-Extra-Antcloud-Aks-Cluster-Id: helm-dxfzuzcn-cluster" \
-H "X-Remote-Extra-Antcloud-Aks-Tenant-Id: DXFZUZCN" \
-H "X-Remote-Extra-Antcloud-Aks-Workspace-Id: helm" \
-H "X-Remote-Group: system:masters" \
-H "X-Remote-User: system:apiserver" \
-H "User-Agent: hyperkube/v1.12.0 (darwin/amd64) kubernetes/31dda1c" \
-H "X-Forwarded-For: ::1" \
-H "Accept-Encoding: gzip" \
-H "Accept: application/vnd.kubernetes.protobuf, */*"  \
https://localhost:443/apis/metrics.k8s.io/v1beta1/namespaces/kube-system/pods


curl -k --cert ~/kubernetes/aks/cert.pem --key ~/kubernetes/aks/key.pem \
-H "X-Remote-Extra-Antcloud-Aks-Cluster-Id: meshsidercar-test" \
-H "X-Remote-Extra-Antcloud-Aks-Tenant-Id: ETOPAFCN" \
-H "X-Remote-Extra-Antcloud-Aks-Workspace-Id: meshsidecartest1" \
-H "X-Remote-Group: system:masters" \
-H "X-Remote-User: system:apiserver" \
-H "User-Agent: hyperkube/v1.12.0 (darwin/amd64) kubernetes/31dda1c" \
-H "X-Forwarded-For: ::1" \
-H "Accept-Encoding: gzip" \
-H "Accept: application/vnd.kubernetes.protobuf, */*"  \
https://localhost:443/apis/metrics.k8s.io/v1beta1/namespaces/kube-system/pods


curl -k --cert ~/kubernetes/aks/cert.pem --key ~/kubernetes/aks/key.pem \
-H "X-Remote-Extra-Antcloud-Aks-Cluster-Id: meshsidercar-test" \
-H "X-Remote-Extra-Antcloud-Aks-Tenant-Id: ETOPAFCN" \
-H "X-Remote-Extra-Antcloud-Aks-Workspace-Id: meshsidecartest1" \
-H "X-Remote-Group: system:masters" \
-H "X-Remote-User: system:apiserver" \
-H "User-Agent: hyperkube/v1.12.0 (darwin/amd64) kubernetes/31dda1c" \
-H "X-Forwarded-For: ::1" \
-H "Accept-Encoding: gzip" \
-H "Accept: application/vnd.kubernetes.protobuf, */*"  \
https://localhost:443/apis/metrics.k8s.io/v1beta1/namespaces/default/pods


curl -k --cert ~/kubernetes/aks/cert.pem --key ~/kubernetes/aks/key.pem \
-H "X-Remote-Extra-Antcloud-Aks-Cluster-Id: meshsidercar-test" \
-H "X-Remote-Extra-Antcloud-Aks-Tenant-Id: ETOPAFCN" \
-H "X-Remote-Extra-Antcloud-Aks-Workspace-Id: meshsidecartest1" \
-H "X-Remote-Group: system:masters" \
-H "X-Remote-User: system:apiserver" \
-H "User-Agent: hyperkube/v1.12.0 (darwin/amd64) kubernetes/31dda1c" \
-H "X-Forwarded-For: ::1" \
-H "Accept-Encoding: gzip" \
-H "Accept: application/vnd.kubernetes.protobuf, */*"  \
https://localhost:443/apis/metrics.k8s.io/v1beta1/namespaces/default/pods/reviews-v2-dc7756457-dq8v9


curl -k --cert ~/kubernetes/aks/cert.pem --key ~/kubernetes/aks/key.pem \
-H "X-Remote-Extra-Antcloud-Aks-Cluster-Id: helmtest" \
-H "X-Remote-Extra-Antcloud-Aks-Tenant-Id: MEUGIVCN" \
-H "X-Remote-Extra-Antcloud-Aks-Workspace-Id: HelmTest" \
-H "X-Remote-Group: system:masters" \
-H "X-Remote-User: system:apiserver" \
-H "User-Agent: hyperkube/v1.12.0 (darwin/amd64) kubernetes/31dda1c" \
-H "X-Forwarded-For: ::1" \
-H "Accept-Encoding: gzip" \
-H "Accept: application/vnd.kubernetes.protobuf, */*"  \
https://localhost:443/apis/metrics.k8s.io/v1beta1/namespaces/default/pods

curl -k -H "X-Remote-Extra-Antcloud-Aks-Cluster-Id: helmtest" \
-H "X-Remote-Extra-Antcloud-Aks-Tenant-Id: MEUGIVCN" \
-H "X-Remote-Extra-Antcloud-Aks-Workspace-Id: HelmTest" \
-H "X-Remote-Group: system:masters" \
-H "X-Remote-User: system:apiserver" \
-H "User-Agent: hyperkube/v1.12.0 (darwin/amd64) kubernetes/31dda1c" \
-H "X-Forwarded-For: ::1" \
-H "Accept-Encoding: gzip" \
-H "Accept: application/vnd.kubernetes.protobuf, */*"  \
https://localhost:6443/apis/metrics.k8s.io/v1beta1/namespaces/default/pods

