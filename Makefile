default: ;
.DEFAULT_GOAL: default

benchmark:
	go run cmd/benchmark/main.go
.PHONY: benchmark

fund:
	go run cmd/fund/main.go
.PHONY: fund

glide_version := v0.13.1
glide_arch := linux-amd64
glide:
	curl -sSLo glide.tar.gz https://github.com/Masterminds/glide/releases/download/$(glide_version)/glide-$(glide_version)-$(glide_arch).tar.gz
	tar -xf ./glide.tar.gz
	mv ./$(glide_arch)/glide ./glide
	rm -rf ./$(glide_arch) ./glide.tar.gz
