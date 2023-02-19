FROM golang:1.16.2 as build
RUN mkdir -p /app
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o=/app/bin/calendar /app/cmd/calendar

FROM alpine:latest as app
COPY --from=build /app/bin/calendar /app/bin/calendar
COPY --from=build /app/configs /app/configs
RUN chmod +x /app/bin/calendar

ENTRYPOINT ["/app/bin/calendar"]