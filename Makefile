.PHONY: migrate-up migrate-down all

ROOT_PATH := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

include main.env

all: clean
	go build -o build ./...

generate-migration:
	migrate create -ext sql -dir migrations $(NAME)

migrate-up:
	migrate -path $(ROOT_PATH)/migrations -database $(PG_URL) up $(step)

migrate-down:
	migrate -path $(ROOT_PATH)/migrations -database $(PG_URL) down $(step)

clean:
	rm -f ./build/*