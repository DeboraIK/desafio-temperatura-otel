FROM golang:1.23 AS build
WORKDIR /app
ARG FOLDER=api_b
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /api cmd/${FOLDER}/main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /
COPY --from=build /api /api
ENTRYPOINT ["/api"]
