dev-build:
	go build -o policy_asserter .
build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o policy_asserter .
dockerize: build
	docker build . --tag $DOCKER_REPO/sslpolicy-asserter:1.2
publish: dockerize
	docker push $DOCKER-REPO/sslpolicy-asserter:1.2
