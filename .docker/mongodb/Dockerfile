#
# MongoDB Dockerfile
#

FROM mongo:latest

# Generate the mongo replication keyFile via openssl
RUN openssl rand -base64 756 > /etc/mongo.replica.key

# Change ownership and permissions of the replication key file
RUN chmod 400 /etc/mongo.replica.key
RUN chown mongodb:mongodb /etc/mongo.replica.key
