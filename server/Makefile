DATABASE_URL := "postgres://postgres:postgres@localhost:54321/backups?sslmode=disable"

migrations:
	goose --dir=migration postgres ${DATABASE_URL} up