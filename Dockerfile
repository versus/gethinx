FROM alpine:edge AS build
RUN apk update
RUN apk upgrade
RUN apk add --update go gcc g++ git
WORKDIR /app
ENV GOPATH /app
RUN mkdir -p  /app/src/gethinx
ADD . /app/src/gethinx
RUN cd /app && go get gethinx # server is name of our application
RUN CGO_ENABLED=1 GOOS=linux go install -a gethinx


FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

EXPOSE 8080
EXPOSE 8545
# Add our application binary
COPY --from=build /app/bin/gethinx /app/gethinx
ADD ./artefacts/templates  /app/templates
ADD ./artefacts/config.toml  /app/config.toml
WORKDIR /app
ENTRYPOINT [ "/app/gethinx", "-c", "/app/config.toml" ]
