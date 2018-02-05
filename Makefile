TEST?=./...

.PHONY: config-install get-deps code-analysis test all

all: get-deps code-analysis test

config-install:
	mkdir -p "$(HOME)/.config/perun"
	cp defaults/main.yaml "$(HOME)/.config/perun/main.yaml"

get-deps:
	go get -t -v .
	go get github.com/stretchr/testify
	go install github.com/stretchr/testify
	go build
	go fmt ...

code-analysis: get-deps
	go tool vet -v .

test: get-deps
	go test -v -cover $(TEST)

