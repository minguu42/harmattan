FROM golang:1.20 AS base
WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

FROM base AS dev
RUN go install github.com/cosmtrek/air@latest
CMD ["air", "-c", ".air.toml"]

FROM base AS build
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/app ./cmd/server

FROM gcr.io/distroless/static-debian11 AS prod
COPY --from=build /go/bin/app /
ENTRYPOINT ["/app"]
