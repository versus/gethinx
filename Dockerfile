FROM ubuntu:xenial
#FROM debian:jessie



EXPOSE 8080
EXPOSE 8545

# Add our application binary
ADD ./artefacts/gethinx-linux-x64 /app/gethinx
ADD ./artefacts/templates  /app/templates
ADD ./artefacts/config.toml  /app/config.toml
RUN chown -R nobody:nogroup /app

USER nobody
WORKDIR /app
ENTRYPOINT [ "/app/gethinx", "-c", "/app/config.toml" ]
