.PHONY: all
all: bookstore

.PHONY: all_build
all_build:
	@printf "Building bookstore service"
	@$(MAKE) bookstore_build
	@printf "\rBuilding bookstore service [DONE]\n"
	@printf "Building migration engine"
	@$(MAKE) migrate_build
	@printf "\rMigration engine build [DONE]\n"

.PHONY: bookstore
bookstore: bookstore_build bookstore_run

.PHONY: bookstore_build
bookstore_build: 
	@go build -o bin/bookstore ./cmd/bookstore

.PHONY: bookstore_run
bookstore_run:
	@./bin/bookstore

.PHONY: migrate_build
migrate_build: 
	@go build -o bin/migrate ./cmd/migrate

.PHONY: clean
clean:
	@rm -rf bin/**/*
