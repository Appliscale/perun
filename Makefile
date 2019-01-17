.PHONY: get-deps code-analysis test all

all: get-deps code-analysis test

get-deps:
	go get -t -v ./...
	go install ./...
	go build
	go fmt ./...

code-analysis: get-deps
	go vet -v ./...

test: get-deps create-mocks
	go test -cover ./...

create-mocks: get-mockgen
	`go env GOPATH`/bin/mockgen -source=./awsapi/cloudformation.go  -destination=./stack/stack_mocks/mock_aws_api.go -package=stack_mocks CloudFormationAPI
	`go env GOPATH`/bin/mockgen -source=./logger/logger.go  -destination=./checkingrequiredfiles/mocks/mock_logger.go -package=mocks LoggerInt

get-mockgen:
	go get github.com/golang/mock/gomock
	go install github.com/golang/mock/mockgen
