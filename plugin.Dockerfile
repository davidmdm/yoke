FROM golang:1.22-alpine

WORKDIR /cmp

COPY go.mod go.sum ./

RUN go mod download

COPY ./cmd/yokecd ./cmd/yokecd
COPY ./internal ./internal

RUN go install ./cmd/yokecd

COPY ./plugin.yaml /home/argocd/cmp-server/config/plugin.yaml

