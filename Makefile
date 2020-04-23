CONTAINER_TAG ?= 1.2
CONTAINER_REPO ?= ct-services
dev-build:
	go build -o policy_asserter .
build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o policy_asserter .
dockerize: build
	docker build . --tag '$(CONTAINER_REPO)'/sslpolicy-asserter:'$(CONTAINER_TAG)'
publish: dockerize
	docker push '$(CONTAINER-REPO)'/sslpolicy-asserter:'$(CONTAINER_TAG)'
