FROM golang:1.22-bullseye

WORKDIR /usr/src/

RUN go install github.com/cosmtrek/air@latest && \
    go install github.com/a-h/templ/cmd/templ@latest

COPY go.mod go.sum ./

RUN go mod download

COPY ./.env ./.env
COPY ./cmd ./cmd
COPY ./config ./config
COPY ./internal ./internal
COPY ./web ./web

RUN templ generate && \
    go mod tidy

CMD ["air", "-c", "config/.air.toml"]
