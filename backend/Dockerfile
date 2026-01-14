FROM golang:1.24-alpine AS build
WORKDIR /src
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/app ./cmd/server
RUN CGO_ENABLED=0 go build -o /out/seed ./cmd/seed

FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/app ./app
COPY --from=build /out/seed ./seed
COPY migrations ./migrations
EXPOSE 8080
ENTRYPOINT ["./app"]