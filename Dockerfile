FROM golang:1.18-alpine3.14 AS builder
ARG TARGETARCH
ENV CGO_ENABLED=0
COPY . /app
WORKDIR /app/cmd
RUN GOARCH=$TARGETARCH go build -o ws-chat .

FROM alpine:3.14
COPY --from=builder /app/cmd/ws-chat /app/cmd/ws-chat
COPY --from=builder /app/static /app/static
WORKDIR /app/cmd/
CMD ./ws-chat
