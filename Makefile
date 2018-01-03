default: ;
.DEFAULT_GOAL: default

horizon_addr := http://localhost:8000

benchmark:
	go run cmd/benchmark/main.go -address $(horizon_addr)
.PHONY: benchmark

fund:
	go run cmd/fund/main.go
.PHONY: fund

cat-horizon-logs:
	docker-compose exec horizon sh -c 'cat /var/log/supervisor/horizon-stderr*'
.PHONY: cat-horizon-logs

tail-horizon-logs:
	docker-compose exec horizon sh -c 'tail -f /var/log/supervisor/horizon-stderr*'
.PHONY: tail-horizon-logs

cat-stellar-core-logs:
	docker-compose exec horizon sh -c 'cat /var/log/supervisor/stellar-core-stdout*'
.PHONY: cat-stellar-core-logs

tail-stellar-core-logs:
	docker-compose exec horizon sh -c 'tail -f /var/log/supervisor/stellar-core-stdout*'
.PHONY: tail-stellar-core-logs

glide_version := v0.13.1
glide_arch := linux-amd64
glide:
	curl -sSLo glide.tar.gz https://github.com/Masterminds/glide/releases/download/$(glide_version)/glide-$(glide_version)-$(glide_arch).tar.gz
	tar -xf ./glide.tar.gz
	mv ./$(glide_arch)/glide ./glide
	rm -rf ./$(glide_arch) ./glide.tar.gz
