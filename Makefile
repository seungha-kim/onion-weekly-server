main:
	go build -o build/onion-weekly onion-weekly.go

migrate-up:
	migrate -path ./migrations -database "postgresql://postgres@localhost/onion_weekly?sslmode=disable" up $(step)

migrate-down:
	migrate -path ./migrations -database "postgresql://postgres@localhost/onion_weekly?sslmode=disable" down $(step)