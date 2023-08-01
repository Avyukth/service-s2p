SHELL 			:= /bin/bash
VERSION 		:= 1.0
KIND_CLUSTER    := subhrajit-starter-cluster
KIND            := kindest/node:v1.27.3

run:
	go run app/services/sales-api/main.go  |  go run app/tooling/logfmt/main.go


all: sales-api

sales-api:
	docker build \
	--no-cache \
	-f zarf/docker/dockerfile.sales-api \
	-t sales-api-amd64:$(VERSION) \
	--build-arg BUILD_REF=$(VERSION) \
	--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
	.


kind-up:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/kind/kind-config.yaml
	kubectl config set-context --current --namespace=sales-system

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)


kind-load:
	cd zarf/k8s/kind/sales-pod; kustomize edit set image sales-api-image=sales-api-amd64:$(VERSION)
	kind load docker-image sales-api-amd64:$(VERSION) --name $(KIND_CLUSTER)

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-apply:
	kustomize build zarf/k8s/kind/sales-pod | kubectl apply -f -

kind-status-sales:
	kubectl get pods -o wide --watch

kind-logs:
	kubectl logs -l app=sales --all-containers=true -f --tail=100

kind-restart:
	kubectl rollout restart deployment sales-pod

kind-update: all kind-load kind-restart

kind-describe:
	kubectl describe pod -l app=sales

kind-update-apply: all kind-load kind-apply


# =================================================================
# Module Vendor Support
tidy:
	go mod tidy
	go mod vendor
