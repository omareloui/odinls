FROM golang:1.24-bullseye

# For the bash commands to work properly
SHELL ["/bin/bash", "-o", "pipefail", "-c"]

WORKDIR /usr/src/

RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash -  && \
    apt-get -y --no-install-recommends install nodejs && \
    npm i -g tailwindcss && \
    \
    go install github.com/air-verse/air@v1.62.0 && \
    go install github.com/a-h/templ/cmd/templ@v0.3.906

COPY go.mod go.sum ./

RUN go mod download

COPY Makefile .
COPY .env .
COPY cmd cmd
COPY config config
COPY internal internal
COPY web web

RUN templ generate && go mod tidy


ENTRYPOINT [ "air" ]
CMD ["-c", "config/.air.toml"]
