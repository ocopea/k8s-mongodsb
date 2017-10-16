# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.8

RUN curl -o /usr/local/bin/swagger -L'#' https://github.com/go-swagger/go-swagger/releases/download/0.10.0/swagger_$(echo `uname`|tr '[:upper:]' '[:lower:]')_amd64
RUN chmod +x /usr/local/bin/swagger

# Copy the local package files to the container's workspace and build.
ADD . /go/src/ocopea/k8s-mongodsb
RUN cd /go/src/ocopea/k8s-mongodsb && swagger generate server -f dsb-swagger.yaml -A k8s-mongodsb
RUN go install ocopea/k8s-mongodsb/cmd/k8s-mongodsb-server

# Copy Mongo utilities
ADD target/mongo/mongodump /go/bin/mongodump
ADD target/mongo/mongorestore /go/bin/mongorestore
ADD target/mongo/mongo /go/bin/mongoshell

# Cleanup source code artifacts
RUN rm -rf /go/src/*

ENTRYPOINT /go/bin/k8s-mongodsb-server

EXPOSE 8000