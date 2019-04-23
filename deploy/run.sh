#!/usr/bin/env bash

# shanghai-finance-pre
./metrics-apiserver \
--lister-kubeconfig=/etc/metrics-apiserver/conf/kubeconfig.yml \
--authentication-kubeconfig=/etc/metrics-apiserver/conf/kubeconfig.yml \
--authorization-kubeconfig=/etc/metrics-apiserver/conf/kubeconfig.yml \
--kube-config=/etc/metrics-apiserver/conf/kubeconfig.yml \
--client-ca-file=/etc/metrics-apiserver/conf/ca.pem \
--requestheader-client-ca-file=/etc/metrics-apiserver/conf/ca.pem \
--secure-port=6443 \
--authentication-skip-lookup=true \
--monitor-server-url=http://10.252.4.12:8341


curl -k --cert /root/metrics-apiserver/conf/cert.pem --key /root/metrics-apiserver/conf/key.pem \
-H "X-Remote-Extra-Antcloud-Aks-Cluster-Id: test" \
-H "X-Remote-Extra-Antcloud-Aks-Tenant-Id: XOVKOVCN" \
-H "X-Remote-Extra-Antcloud-Aks-Workspace-Id: ShanghaiTest2" \
-H "X-Remote-Group: system:masters" \
-H "X-Remote-User: system:apiserver" \
-H "User-Agent: hyperkube/v1.12.0 (darwin/amd64) kubernetes/31dda1c" \
-H "X-Forwarded-For: ::1" \
-H "Accept-Encoding: gzip" \
-H "Accept: application/vnd.kubernetes.protobuf, */*"  \
https://localhost:443/apis/metrics.k8s.io/v1beta1/namespaces/default/pods