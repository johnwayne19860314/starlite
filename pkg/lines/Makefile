check: generate vendor test lint

generate-pen:
	cd tools/pen && make

generate:
	cd example/petstore/ && make
	cd example/pkg/ && make pet-store

test:
	go test ./... -p 1

lint:
	golangci-lint run --timeout 1200s --verbose ./...

test-coverage:
	go test ./... --coverprofile testcover
	go tool cover -html=testcover -o testcover.html
	open testcover.html

vendor:
	go mod vendor

build-tools:
	docker build -f tools.Dockerfile \
		-t china-dots-system-docker-local.arf.xxx.cn/lines/tools:latest .

run-tools:
	docker run --rm china-dots-system-docker-local.arf.xxx.cn/lines/tools:latest $(ARGS)
