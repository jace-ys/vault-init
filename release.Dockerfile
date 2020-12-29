FROM scratch
COPY vault-init /bin/vault-init
ENTRYPOINT ["vault-init"]
