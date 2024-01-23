
@PHONY: run debug build bin_dir

all: debug

bin_dir:
	mkdir -p bin

build: cmd/main.go | bin_dir
	go build -o bin/main cmd/backend/main.go

run:  | bin_dir
	@DEBUG=false LOG_LEVEL=2 go run cmd/backend/main.go

debug: | bin_dir
	@DEBUG=true go run cmd/backend/main.go
