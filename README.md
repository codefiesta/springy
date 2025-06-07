
# Springy
[![Build](https://github.com/codefiesta/springy/workflows/Go/badge.svg)](https://github.com/codefiesta/springy/actions?query=workflow:Go)

Springy is an open source real-time database backed by [MongoDB](https://github.com/mongodb/mongo) that allows 
developers to publish and subscribe to real time events.

## Getting Started
If you want to build Springy you have two options:

##### You have a working [Go environment] 

```
git clone https://github.com/codefiesta/springy
cd springy
go build go.springy.io
go run go.springy.io
```

##### You have a working [Docker environment] 

```
git clone https://github.com/codefiesta/springy
cd springy
docker-compose up -d --build
```