#
# MongoDB Dockerfile
#

FROM mongo:latest

LABEL maintainer="Springy.io <springy@spring.io>"

# Copy the mongo init script which configures the db
# See https://stackoverflow.com/questions/42912755/how-to-create-a-db-for-mongodb-container-on-start-up
ADD mongo-init.js /docker-entrypoint-initdb.d/

# Define mountable directories.
VOLUME ["/data/db"]

# Define working directory.
WORKDIR /data

# Define default command.
CMD ["mongod"]

# Expose ports.
#   - 27017: process
#   - 28017: http
EXPOSE 27017
EXPOSE 28017