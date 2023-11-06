.PHONY: lint mocks tests

lint:
	golangci-lint run

mocks:
	mockery --case snake --dir ./repositories --all --output ./mocks/repositories

tests:
	go test -v -cover -race -timeout 300s -count=1 ./...

build:
	@cd cmd/${cmd} && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${service_name} .



