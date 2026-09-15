FROM golang:1.26-alpine AS builder
RUN apk add --no-cache ca-certificates git
WORKDIR /src
RUN go install github.com/luno/luno-mcp/cmd/server@v0.6.3
COPY go.mod main.go ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/headortails-mcp .

FROM alpine:3.22
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/bin/server /usr/local/bin/luno-mcp
COPY --from=builder /out/headortails-mcp /usr/local/bin/headortails-mcp
ENV LUNO_API_DOMAIN=api.luno.com
ENV ALLOW_WRITE_OPERATIONS=true
EXPOSE 8080
USER nobody
ENTRYPOINT ["/usr/local/bin/headortails-mcp"]
