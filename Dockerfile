FROM golang:1.20 AS base
WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    go mod download

FROM base AS dev
RUN go install github.com/cosmtrek/air@latest
CMD ["air", "-c", ".air.toml"]

FROM base AS build
ARG APP_VERSION="v0.0.0+unknown"
ARG APP_REVISION=""
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 go build -ldflags "-X main.version=$APP_VERSION -X main.revision=$APP_REVISION" -o /go/bin/app ./cmd/server

FROM gcr.io/distroless/static-debian11 AS prod
COPY --from=build /go/bin/app /
ENTRYPOINT ["/app"]
