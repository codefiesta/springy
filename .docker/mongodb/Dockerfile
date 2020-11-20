#
# MongoDB Dockerfile
#

FROM mongo:latest

LABEL maintainer="Springy.io <springy@spring.io>"

# Copy the mongo init script which configures the db
# See https://zgadzaj.com/development/docker/docker-compose/turning-standalone-mongodb-server-into-a-replica-set-with-docker-compose#env-dist
ADD mongo-init.sh /docker-entrypoint-initdb.d/

# Define mountable directories.
VOLUME ["/data/db"]

# Define working directory.
WORKDIR /data

# Define default command.
#TODO: Find out how to grab the $MONGO_REPLICA_SET environment variable
CMD ["mongod", "--replSet", "springyRs"]

# Expose ports.
#   - 27017: process
#   - 28017: http
EXPOSE 27017
EXPOSE 28017