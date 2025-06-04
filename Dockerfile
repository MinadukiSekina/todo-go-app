FROM golang:latest

ARG UID=20000
ARG GID=20000
RUN groupadd -g $GID dev \
    && useradd -m -u $UID -g $GID dev
# bashに変更
RUN chsh -s /bin/bash dev
USER dev

WORKDIR /app
COPY --chown=dev:dev . .

RUN go mod tidy

RUN go install -v github.com/air-verse/air@latest
RUN go install -v github.com/go-delve/delve/cmd/dlv@latest

