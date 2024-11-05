include .env

dev:
	@go run cmd/tex/main.go
db:
	@echo "Initializing texApi database..."
	@PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d postgres \
		-f ./schemas/0.4.1_create_landing.sql
	@PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d db_tex \
		-f ./schemas/0.4.2_insert_landing.sql
	@PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d db_tex \
		-f ./schemas/0.5.1_create_core.sql
	@PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d db_tex \
			-f ./schemas/0.5.2_logisticops.sql
	@echo "Has been successfully created"
build:
	@echo "Building the app, please wait..."
	@go build -o ./bin/texApi cmd/tex/main.go
	@echo "Done."
build-cross:
	@echo "Bulding for windows, linux and macos (darwin m2), please wait..."
	@GOOS=linux GOARCH=amd64 go build -o ./bin/texApi-linux cmd/tex/main.go
	@GOOS=darwin GOARCH=arm64 go build -o ./bin/texApi-macos cmd/tex/main.go
	@GOOS=windows GOARCH=amd64 go build -o ./bin/texApi-windows cmd/tex/main.go
	@echo "Done."
