FROM golang:1.16

RUN mkdir -p /integration_tests
WORKDIR /integration_tests

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
CMD ["go", "test"]
