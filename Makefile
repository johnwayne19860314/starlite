check: vendor test lint

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
