FROM alpine:latest

COPY --chmod=0755 --chown=root:root ./swarm-key-maker /usr/local/bin/swarm-key-maker
ENTRYPOINT ["/usr/local/bin/swarm-key-maker"]
