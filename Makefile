.PHONY: config-install get-deps code-analysis test all

all: get-deps code-analysis test

config-install:
	mkdir -p "$(HOME)/.config/perun"
	cp defaults/main.yaml "$(HOME)/.config/perun/main.yaml"
	mkdir -p "$(HOME)/perun"
	cp stack/defaultstackpolicy/blocked.json "$(HOME)/perun/blocked.json"
	cp stack/defaultstackpolicy/unblocked.json "$(HOME)/perun/unblocked.json"

get-deps:
	go get -t -v ./...
	go install ./...
	go build
	go fmt ./...

code-analysis: get-deps
	go vet -v ./...

test: get-deps
	go test -v -cover ./...
