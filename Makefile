.PHONY: migrate-up migrate-down

include $(ENV).env

migrate-up:
	migrate -path ./migrations -database $(PG_URL) up $(step)

migrate-down:
	migrate -path ./migrations -database $(PG_URL) down $(step)
