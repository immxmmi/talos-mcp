FROM cgr.dev/chainguard/wolfi-base:latest AS build

RUN apk add --no-cache go git

WORKDIR /app
COPY src/ .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /talos-mcp .

FROM cgr.dev/chainguard/static:latest
COPY --from=build /talos-mcp /usr/bin/talos-mcp
ENTRYPOINT ["/usr/bin/talos-mcp"]
