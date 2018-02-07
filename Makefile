default: ;
.DEFAULT_GOAL: default

vendor:
	./glide install

build:
	go build -o loadtest cmd/loadtest/*.go
	go build -o create cmd/create/*.go
	go build -o merge cmd/merge/*.go
.PHONY: loadtest

glide_version := v0.13.1
glide_arch := linux-amd64
glide:
	curl -sSLo glide.tar.gz https://github.com/Masterminds/glide/releases/download/$(glide_version)/glide-$(glide_version)-$(glide_arch).tar.gz
	tar -xf ./glide.tar.gz
	mv ./$(glide_arch)/glide ./glide
	rm -rf ./$(glide_arch) ./glide.tar.gz

clean:
	rm -f loadtest
	rm -f create
	rm -f merge
	rm -f glide
	rm -rf vendor
.PHONY: clean
