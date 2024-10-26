
@PHONY: run debug build bin_dir create_db

all: debug

install-go:
	go mod download

install-js:
	(cd ui && pnpm install)

install: install-js install-go

format-go:
	goimports -w $(shell find . -type f -name '*.go' -not -path "./graph/*")
	gofmt -w $(shell find . -type f -name '*.go' -not -path "./graph/*")

format-js:
	(cd ui && pnpm format)

format: format-js format-go

test-go:
	go test --race -coverprofile=coverage.txt -covermode=atomic ./...


build-js:
	(cd ui && pnpm build)

build-go:
	go build ./cmd/backend/main.go

build: build-js build-go

create_db:
	docker run \
		--detach \
		--name "firepit-mariadb" \
		--env MARIADB_ROOT_PASSWORD="root" \
		--env MARIADB_DATABASE="firepit-mariadb" \
		-p 127.0.0.1:3306:3306 \
		mariadb:latest
