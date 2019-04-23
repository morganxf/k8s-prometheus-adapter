############################
# STEP 1 build executable binary
############################
FROM golang:1.11 AS builder

RUN apt-get update \
    && apt-get install -y vim-tiny

ARG VERSION
WORKDIR /go/src/github.com/directxman12/k8s-prometheus-adapter
ADD . .
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -tags=multitenancy -ldflags "-X main.version=${VERSION}" -o metrics-apiserver github.com/directxman12/k8s-prometheus-adapter/cmd/metrics-apiserver

############################
# STEP 2 build a small image
############################
FROM debian:stretch

RUN apt-get update \
    && apt-get install -y vim-tiny procps --no-install-recommends

ENV HOME=/etc/metrics-apiserver

WORKDIR /etc/metrics-apiserver
# Copy our static executable and configuration.
COPY --from=builder /go/src/github.com/directxman12/k8s-prometheus-adapter/metrics-apiserver ./
COPY scripts/entrypoint.sh ./
RUN chmod +x /etc/metrics-apiserver/entrypoint.sh

ENTRYPOINT ["/etc/metrics-apiserver/entrypoint.sh"]
CMD [ "/etc/metrics-apiserver/metrics-apiserver", \
      "--lister-kubeconfig=/etc/metrics-apiserver/conf/kubeconfig.yml", \
      "--authentication-kubeconfig=/etc/metrics-apiserver/conf/kubeconfig.yml", \
      "--authorization-kubeconfig=/etc/metrics-apiserver/conf/kubeconfig.yml", \
      "--kube-config=/etc/metrics-apiserver/conf/kubeconfig.yml", \
      "--client-ca-file=/etc/metrics-apiserver/conf/ca.pem", \
      "--requestheader-client-ca-file=/etc/metrics-apiserver/conf/ca.pem", \
      "--secure-port=6443", \
      "--authentication-skip-lookup=true" ]
