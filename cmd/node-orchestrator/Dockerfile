FROM alpine:latest

COPY --chmod=0755 --chown=root:root ./node-orchestrator /usr/local/bin/node-orchestrator
ENTRYPOINT ["/usr/local/bin/node-orchestrator"]
