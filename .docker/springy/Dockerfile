#
# Spring Microservice Dockerfile
#

FROM golang:1.18.4-alpine

LABEL maintainer="Springy.io <springy@spring.io>"

# Alpine container is limited so add bash, git, openssh
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

# Define working directory.
WORKDIR /springy

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Compile the Go app into an executable
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["./main"]
