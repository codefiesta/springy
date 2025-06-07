# Springy Docker Setup Script
brew install go openssl
brew install docker --cask
brew install --cask mongodb-compass

# Run the docker command
docker-compose up -d --build
