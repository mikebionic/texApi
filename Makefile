include .env

dev:
	@go run cmd/tex/main.go

db:
	@echo "Initializing texApi database..."
	@bash ./scripts/re-init-db.sh
#	@PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d postgres \
#    -f ./schemas/0.0.5_drop_db.sql

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

upload-dir:
	@mkdir -p $(UPLOAD_PATH) || (echo "Error: Failed to create directory $(UPLOAD_PATH)" && exit 1)
	@echo "Directory $(UPLOAD_PATH) created"

init-sys:
	@mkdir -p ~/tex_backend/app/
	@cp -r ~/tex_backend/texApi/assets/ ~/tex_backend/app/
	@cp ~/tex_backend/texApi/.env.example ~/tex_backend/app/.env
	@sudo cp scripts/texApi.service /etc/systemd/system/texApi.service