.PHONY: migrate-up migrate-down all

include main.env

all: clean
	go build -o build ./...

generate-migration:
	migrate create -ext sql -dir migrations $(NAME)

migrate-up:
	migrate -path ./migrations -database $(PG_URL) up $(step)

migrate-down:
	migrate -path ./migrations -database $(PG_URL) down $(step)

clean:
	rm -f ./build/*