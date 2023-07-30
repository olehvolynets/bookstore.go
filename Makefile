.PHONY: all
all: build_server run_server

.PHONY: build_server
build_server: 
	@go build -o bin/bookstore ./cmd/server

.PHONY: run_server
run_server:
	@./bin/bookstore s

.PHONY: clean
clean:
	@rm -rf bin/**/*
