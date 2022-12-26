FROM golang:alpine AS builder


WORKDIR /
ADD go.mod .
COPY . .
RUN go build -o bin/spotify-service.exe -ldflags="-s -w"

FROM alpine

ARG PNKL_VAULT_TOKEN
ARG GO_ENVIRONMENT
ARG LOG_LEVEL

ENV PNKL_VAULT_TOKEN=$PNKL_VAULT_TOKEN
ENV GO_ENVIRONMENT=$GO_ENVIRONMENT
ENV LOG_LEVEL=$LOG_LEVEL

RUN apk update
WORKDIR /
COPY --from=builder /bin .
EXPOSE 8080
CMD ["./spotify-service.exe"]