# How to



## create migration file

`go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate/`

`migrate create -ext sql -dir migrations <migration_name>`

`migrate -path ./migrations -database "postgresql://postgres@localhost/onion_weekly?sslmode=disable" up`