SHELL := /bin/bash

# ==============================================================================
# Testing running system

# For testing a simple query on the system. Don't forget to `make seed` first.
# curl --user "admin@example.com:sup3rS3cr3tGolang" http://localhost:3000/users/token/8ea09532-9245-8623-923d-3201212966b1
# export TOKEN="COPY TOKEN STRING FROM LAST CALL"
# curl -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/users/1/2

# For testing load on the service.
# hey -m GET -c 100 -n 10000 -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/users/1/2
# expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,mem:memstats.Alloc"

# ==============================================================================
# Building containers

all: sales-api

sales-api:
	docker build \
		-f zarf/docker/dockerfile.sales-api \
		-t sales-api-amd64:1.0 \
		--build-arg VCS_REF=`git rev-parse HEAD` \
		--build-arg BUILD_DATE=`date -u +”%Y-%m-%dT%H:%M:%SZ”` \
		.

# ==============================================================================
# Running from within k8s/dev

kind-up:
	kind create cluster --image kindest/node:v1.20.2 --name neo-starter-cluster --config zarf/k8s/dev/kind-config.yaml

kind-down:
	kind delete cluster --name neo-starter-cluster

kind-load:
	kind load docker-image sales-api-amd64:1.0 --name neo-starter-cluster

kind-services:
	kustomize build zarf/k8s/dev | kubectl apply -f -

kind-status:
	kubectl get nodes
	kubectl get pods --watch

kind-status-full:
	kubectl describe pod -lapp=sales-api

kind-logs:
	kubectl logs -lapp=sales-api --all-containers=true -f

kind-sales-api: sales-api
	kind load docker-image sales-api-amd64:1.0 --name neo-starter-cluster
	kubectl delete pods -lapp=sales-api

# ==============================================================================

run:
	go run app/sales-api/main.go

runa:
	go run app/admin/main.go

test:
	go test -v ./... -count=1
	staticcheck ./...

tidy:
	go mod tidy
	go mod vendor