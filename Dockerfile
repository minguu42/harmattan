FROM golang:1.21 AS base
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
ARG API_VERSION
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 go build \
      -ldflags "-s -w -X github.com/minguu42/opepe/pkg/handler.version=$API_VERSION" \
      -trimpath \
      -o /go/bin/server \
      ./cmd/server

FROM gcr.io/distroless/static-debian11 AS prod
COPY --from=build /go/bin/server /
ENTRYPOINT ["/server"]
