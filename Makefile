
GO_TEST_TIMEOUT="900s"
GO_TEST_ADDITIONAL_FLAGS :=

GO_JUNIT_REPORT=$(shell which $(GOBIN)/go-junit-report)

####################################
#	Dependencies
####################################
deps: go-junit-report

go-junit-report:
	go get -modfile=tools.mod github.com/jstemmer/go-junit-report@v0.0.0-20190106144839-af01ea7f8024

.PHONY: run
run: deps
	go run --race main.go

.PHONY: integration-test
integration-test: GO_TEST_ADDITIONAL_FLAGS=-tags=integration
integration-test: test

.PHONY: test
test: deps
	go test $(GO_TEST_ADDITIONAL_FLAGS) -covermode=atomic -count=1 -timeout ${GO_TEST_TIMEOUT} -coverprofile=./coverage.txt ./... 2>&1 | tee ./test.txt
	cat test.txt | $(GO_JUNIT_REPORT) > ./report.xml
	go tool cover -func=./coverage.txt

.PHONY: docker-build
docker-build:
	docker build -f ./Dockerfile -t trackmyfish .

.PHONY: docker-tag
docker-tag:
	docker tag trackmyfish simondrake/trackmyfish:v1alpha1

.PHONY: docker-push
docker-push:
	docker push simondrake/trackmyfish:v1alpha1

.PHONY: docker-all
docker-all: docker-build docker-tag docker-push
	@echo "================================================"
	@echo "Docker image has been built, tagged and pushed"
	@echo "================================================"

.PHONY: gosec
run-gosec:
	gosec ./...

.PHONY: lint
lint:
	golangci-lint run ./...
