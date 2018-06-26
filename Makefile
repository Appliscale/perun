.PHONY: config-install get-deps code-analysis test all

all: get-deps code-analysis test

config-install:
	mkdir -p "$(HOME)/.config/perun/stack-policies"
	cp defaults/main.yaml "$(HOME)/.config/perun/main.yaml"
	cp defaults/blocked.json "$(HOME)/.config/perun/stack-policies/blocked.json"
	cp defaults/unblocked.json "$(HOME)/.config/perun/stack-policies/unblocked.json"

get-deps:
	go get -t -v ./...
	go install ./...
	go build
	go fmt ./...

code-analysis: get-deps
	go vet -v ./...

test: get-deps
	go test -v -cover ./...
