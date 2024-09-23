include .env

dev:
	@go run main.go
db:
	@echo "Initializing texApi database..."
	@psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d postgres \
		-f ./database/init/init.sql
	@psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d texApi \
		-f ./database/init/create.sql
	@psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d texApi \
		-f ./database/init/insert.sql
	@psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d texApi \
		-f ./database/init/functions.sql
	@echo "Has been successfully created"
build:
	@echo "Building the app, please wait..."
	@go build -o ./bin/texApi cmd/tex/main.go
	@echo "Done."
build-cross:
	@echo "Bulding for windows, linux and macos (darwin m2), please wait..."
	@GOOS=linux GOARCH=amd64 go build -o ./bin/texApi-linux main.go
	@GOOS=darwin GOARCH=arm64 go build -o ./bin/texApi-macos main.go
	@GOOS=windows GOARCH=amd64 go build -o ./bin/texApi-windows main.go
	@echo "Done."
