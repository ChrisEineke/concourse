FROM alpine:latest

COPY --chmod=0755 --chown=root:root ./deployer /usr/local/bin/deployer
ENTRYPOINT ["/usr/local/bin/deployer"]
