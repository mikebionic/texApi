include .env

dev:
	@go run main.go
db:
	@echo "Initializing texApi database..."
	@psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d postgres \
		-f ./schamas/0.4.1_create_landing.sql
	@psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d texApi \
		-f ./schamas/0.4.2_insert_landing.sql
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
