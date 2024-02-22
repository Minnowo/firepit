
@PHONY: run debug build bin_dir create_db

all: debug

bin_dir:
	mkdir -p bin

build: cmd/main.go | bin_dir
	go build -o bin/main cmd/backend/main.go

run:  | bin_dir
	@DEBUG=false LOG_LEVEL=2 go run cmd/backend/main.go

debug: | bin_dir
	@DEBUG=true go run cmd/backend/main.go


create_db:
	docker run \
		--detach \
		--name "firepit-mariadb" \
		--env MARIADB_ROOT_PASSWORD="root" \
		--env MARIADB_DATABASE="firepit-mariadb" \
		-p 127.0.0.1:3306:3306 \
		mariadb:latest
