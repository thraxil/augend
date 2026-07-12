# syntax=docker/dockerfile:1

FROM golang:1.25 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /augend .

FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=build /augend /app/augend
COPY templates /app/templates
COPY media /app/media
EXPOSE 8080
ENTRYPOINT ["/app/augend"]
