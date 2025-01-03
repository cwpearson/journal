FROM golang:1.23.4-bookworm as builder
ARG GIT_SHA="<not provided>"


RUN apt-get update && apt-get upgrade \
 && rm -rf /var/lib/apt/lists/*

ADD config /src/config
ADD database /src/database
Add entries /src/entries
ADD handlers /src/handlers
ADD models /src/models
ADD ollama /src/ollama
Add tags /src/tags
ADD *.go /src/.
ADD go.mod /src/.

RUN cd /src && go mod tidy
RUN cd /src && go build -ldflags "-X ytdlp-site/config.gitSHA=${GIT_SHA} -X ytdlp-site/config.buildDate=$(date +%Y-%m-%d)" -o server *.go

FROM debian:bookworm-slim

RUN apt-get update \
 && apt-get install -y --no-install-recommends --no-install-suggests \
   ca-certificates \
 && rm -rf /var/lib/apt/lists/*

COPY --from=0 /src/server /opt/server
ADD templates /opt/templates
ADD static /opt/static

WORKDIR /opt
CMD ["/opt/server"]