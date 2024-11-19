#!/bin/bash

cd ~/tex_backend/texApi || { echo "Directory not found"; exit 1; }
git pull origin main
sudo systemctl stop texApi.service
make build
make db
cp ~/tex_backend/texApi/bin/texApi ~/tex_backend/app/texApi
sudo systemctl start texApi.service

echo "Update completed successfully."


