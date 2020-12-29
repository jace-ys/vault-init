# Encryption Backends

Encryption backends are used to encrypt the root tokens and unseal keys generated during the Vault initialization process. The same backend should always be used for each initialization of Vault, as each backend can only decrypt data encrypted by itself. These backends are to be specified via the `vault-init` CLI, using the `--encryption` flag followed by its name (eg. `--encryption=local`).

The implementation of these encryption backends can be found under [`pkg/encryption`](../pkg/encryption).

## Local (`local`)

The `local` encryption backend uses the [AES-GCM encryption algorithm](https://www.cryptosys.net/pki/manpki/pki_aesgcmauthencryption.html) as implemented in the [`crpyto`](https://golang.org/pkg/crypto/) package found in Go's standard library.

#### Configuration

- `--local-encryption-secret-key`: The 32-byte secret key to use for encrypting root tokens and unseal keys.

You can use [OpenSSL](https://www.openssl.org/) to generate a secure 32-byte secret key:

```shell
openssl rand -base64 24
```
