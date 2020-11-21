#!/bin/bash
# Initialization script for preparing a new MongoDB instance.
# This script will perform a couple of duties on a fresh install of a mongo instance:
#   1) Initiate replication to enable mongo change streams which drive most events.
#   2) Change the mongo admin default username & password
#   3) Create an application user / password with read+write priviledges
#   4) Finally, create a collection inside our default database
#
# See: https://stackoverflow.com/questions/42912755/how-to-create-a-db-for-mongodb-container-on-start-up
set -Eeuo pipefail

echo "🐳️ Initializing MongoDB .... 🐳"

mongo <<EOF
  var admin = db.getSiblingDB("admin");
  admin.auth("$MONGO_INITDB_ROOT_USERNAME", "$MONGO_INITDB_ROOT_PASSWORD");
  db.createUser({
    user: "$MONGO_INITDB_USERNAME",
    pwd: "$MONGO_INITDB_PASSWORD",
    roles: ["readWrite"]
  });
EOF

