#!/bin/bash
# Initialization script for preparing a new MongoDB instance.
# This script will perform a couple of duties on a fresh install of a mongo instance:
#   1) Change the mongo admin default username & password
#   2) Create an application user / password with read+write priviledges
#   3) Finally, create a collection inside our default database
#
# Unfortunately, we can't initiate replication to enable mongo change streams which drive most events
# which is why we use the healthcheck trick inside the docker compose file.
#
# See: https://stackoverflow.com/questions/42912755/how-to-create-a-db-for-mongodb-container-on-start-up
# See: https://zgadzaj.com/development/docker/docker-compose/turning-standalone-mongodb-server-into-a-replica-set-with-docker-compose
set -e;

# a default non-root role
echo "🐳🐳🐳🐳🐳🐳🐳🐳🐳🐳🐳🐳 [Initializing MongoDB] 🐳🐳🐳🐳🐳🐳🐳🐳🐳🐳🐳🐳"
MONGO_USER_ROLE="dbOwner"

# Create a default user with the readWrite role in the $MONGO_INITDB_DATABASE database.
mongosh <<-EOJS
  print("🌱🌱🌱🌱🌱🌱🌱🌱🌱🌱🌱🌱 [Seeding MongoDB] 🌱🌱🌱🌱🌱🌱🌱🌱🌱🌱🌱🌱");
  use admin;
  rs.status();
  db.auth("$MONGO_INITDB_ROOT_USERNAME", "$MONGO_INITDB_ROOT_PASSWORD");
  use $MONGO_INITDB_DATABASE;
  db.createCollection("$MONGO_COLLECTION");
	db.createUser({
		user: "$MONGO_INITDB_USERNAME",
		pwd: "$MONGO_INITDB_PASSWORD",
    roles: [ { role: "$MONGO_USER_ROLE", db: "$MONGO_INITDB_DATABASE" } ],
	});
  db.getUser("$MONGO_INITDB_USERNAME");
EOJS

echo "✅✅✅✅✅✅✅✅✅✅✅✅ [MongoDB Initialized] ✅✅✅✅✅✅✅✅✅✅✅✅"
