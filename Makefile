default: ;
.DEFAULT_GOAL: default

dep:
	curl -Lo ./dep https://github.com/golang/dep/releases/download/v0.5.4/dep-linux-amd64
	chmod +x ./dep

vendor: dep
	./dep ensure -v

build: vendor
	go build -o resources/loadtest cmd/loadtest/*.go
	go build -o resources/create cmd/create/*.go
	go build -o resources/merge cmd/merge/*.go
	go build -o resources/whitelist cmd/whitelist/*.go
.PHONY: build

clean: clean-bin clean-vendor
.PHONY: clean

clean-bin:
	rm -f loadtest
	rm -f create
	rm -f merge
.PHONY: clean-bin

clean-vendor:
	rm -rf vendor
.PHONY: clean-vendor
