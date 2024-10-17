#!/bin/bash

# Navigate to the texApi directory
cd ~/tex_backend/texApi || { echo "Directory not found"; exit 1; }

# Pull the latest code from the repository
git pull origin main

# Stop the service
sudo systemctl stop texApp_service.service

# Build the app
make build

# Initialize the database
make db

# Copy the built binary to the service directory
cp ~/tex_backend/texApi/bin/texApi ~/tex_backend/app/texApp_service/tex

# Start the service
sudo systemctl start texApp_service.service

echo "Update completed successfully."


