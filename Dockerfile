FROM golang:1.20 AS base
WORKDIR /go/src/api

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    go mod download

FROM base AS local
RUN go install github.com/cosmtrek/air@latest
CMD ["air", "-c", ".air.toml"]

FROM base AS build
ARG API_VERSION="v0.0.0+unknown"
ARG API_REVISION=""
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 go build \
      -ldflags "-X github.com/minguu42/mtasks/pkg/handler.version=$API_VERSION -X github.com/minguu42/mtasks/pkg/handler.revision=$API_REVISION" \
      -o /go/bin/api \
      ./cmd/server

FROM gcr.io/distroless/static-debian11 AS prod
COPY --from=build /go/bin/api /
ENTRYPOINT ["/api"]
