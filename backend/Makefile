.PHONY: build
build:
	go build -o bin/backend cmd/backend/main.go cmd/backend/config.go

test:
	go test ./...

test-coverage:
	@mkdir -p coverage
	go test -coverprofile=coverage/profile.out ./... && \
		go tool cover -html=coverage/profile.out -o coverage/coverage.html
