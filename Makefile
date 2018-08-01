.PHONY: config-install get-deps code-analysis test all

all: get-deps code-analysis test

config-install:
	mkdir -p "$(HOME)/.config/perun/stack-policies"
	cp defaults/main.yaml "$(HOME)/.config/perun/main.yaml"
	cp defaults/style.yaml "$(HOME)/.config/perun/style.yaml"
	cp defaults/specification_inconsistency.yaml "$(HOME)/.config/perun/specification_inconsistency.yaml"
	cp defaults/blocked.json "$(HOME)/.config/perun/stack-policies/blocked.json"
	cp defaults/unblocked.json "$(HOME)/.config/perun/stack-policies/unblocked.json"

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
	GOPATH=`go env GOPATH` ; $(GOPATH)/bin/mockgen -source=./awsapi/cloudformation.go  -destination=./stack/mocks/mock_aws_api.go -package=mocks CloudFormationAPI

get-mockgen:
	go get github.com/golang/mock/gomock
	go install github.com/golang/mock/mockgen
