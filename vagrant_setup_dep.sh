#!/usr/bin/env bash

set -e

# Update package lists
sudo apt-get update -y

# Install required packages
sudo apt-get install -y \
  apt-transport-https \
  ca-certificates \
  curl \
  gnupg \
  lsb-release \
  make

# Add Docker’s official GPG key
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor \
  -o /usr/share/keyrings/docker-archive-keyring.gpg

# Set up the stable repository
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] \
  https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) \
  stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Install Docker Engine
sudo apt-get update -y
sudo apt-get install -y docker-ce docker-ce-cli containerd.io

# Install Docker Compose plugin (official CLI plugin approach)
sudo apt-get install -y docker-compose-plugin

# Optionally add the 'vagrant' user to the 'docker' group for non-root docker usage
sudo usermod -aG docker vagrant

echo "Docker and Docker Compose installed successfully!"
