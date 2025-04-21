FROM golang:1.23.3

RUN go version
ENV GOPATH=/

COPY ./ ./

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh

# build go api
RUN go mod download
RUN go build -o todo-api ./cmd/main.go

CMD ["./todo-api"]