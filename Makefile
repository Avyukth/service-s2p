SHELL := /bin/bash
VERSION := 1.0
KIND_CLUSTER    := subhrajit-starter-cluster
KIND            := kindest/node:v1.27.3

run:
	go run main.go


all: service

service:
	docker build \
	--no-cache \
	-f zarf/docker/dockerfile \
	-t service-amd64:$(VERSION) \
	--build-arg BUILD_REF=$(VERSION) \
	--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
	.


kind-up:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/kind/kind-config.yaml
	kubectl config set-context --current --namespace=service-system

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)


kind-load:
	kind load docker-image service-amd64:$(VERSION) --name $(KIND_CLUSTER)

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-apply:
	cat zarf/k8s/base/service-pod/base-service.yaml | kubectl apply -f -

kind-status-service:
	kubectl get pods -o wide --watch --namespace=service-system

kind-logs:
	kubectl logs -l app=service --all-containers=true -f --tail=100 --namespace=service-system

kind-restart:
	kubectl rollout restart deployment service-pod --namespace=service-system
