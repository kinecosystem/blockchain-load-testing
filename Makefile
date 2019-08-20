default: ;
.DEFAULT_GOAL: default

vendor:
	./dep ensure -v

build:
	go build -o loadtest cmd/loadtest/*.go
	go build -o create cmd/create/*.go
	go build -o merge cmd/merge/*.go
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
