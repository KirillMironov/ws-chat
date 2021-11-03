FROM golang:1.17-alpine3.14 AS builder
ARG TARGETARCH
ENV CGO_ENABLED=0
ADD . /app
WORKDIR /app/cmd
RUN GOARCH=$TARGETARCH go build -o ws-chat .

FROM alpine:3.14
COPY --from=builder /app/cmd/ws-chat ./cmd/ws-chat
COPY --from=builder /app/static ./static
CMD /cmd/ws-chat
