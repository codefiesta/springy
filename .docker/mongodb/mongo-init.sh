#!/bin/bash
# Initialization script for preparing a new MongoDB instance.
# See: https://stackoverflow.com/questions/42912755/how-to-create-a-db-for-mongodb-container-on-start-up

echo "⭐️⭐️⭐️⭐️⭐️ $MONGO_REPLICA_SET, $MONGO_INITDB_DATABASE  ⭐️⭐️⭐️⭐️⭐️⭐️⭐️"
mongo -- "$MONGO_INITDB_DATABASE" <<EOF
  rs.initiate({_id: "$MONGO_REPLICA_SET", members: [{ _id: 0, host: "127.0.0.1:27017" }] });
  var admin = db.getSiblingDB("admin");
  admin.auth("$MONGO_INITDB_ROOT_USERNAME", "$MONGO_INITDB_ROOT_PASSWORD");
  db.createUser({user: "$MONGO_INITDB_USERNAME", pwd: "$MONGO_INITDB_PASSWORD", roles: ["readWrite"]});
  db = db.getSiblingDB("$MONGO_INITDB_DATABASE")
  db.createCollection("users")
EOF