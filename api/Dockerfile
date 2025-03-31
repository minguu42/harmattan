FROM golang:1.24.0 AS local
WORKDIR /go/src/myapp

RUN --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    go mod download

COPY . .

CMD ["bash", "-c", "CGO_ENABLED=0 go build -ldflags \"-s -w\" -trimpath -o /go/bin/myapp ./api && /go/bin/myapp"]

FROM golang:1.24.0 AS build
WORKDIR /go/src/myapp

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    go mod download

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=.,target=. \
    CGO_ENABLED=0 go build \
      -ldflags "-s -w" \
      -trimpath \
      -o /go/bin/myapp \
      ./api

FROM gcr.io/distroless/static-debian12:nonroot AS prod
COPY --chown=nonroot:nonroot --from=build /go/bin/myapp /
ENTRYPOINT ["/myapp"]
