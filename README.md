[![ci-badge]][ci-workflow] [![release-badge]][release-workflow]

[ci-badge]: https://github.com/jace-ys/vault-init/workflows/ci/badge.svg
[ci-workflow]: https://github.com/jace-ys/vault-init/actions?query=workflow%3Aci
[release-badge]: https://github.com/jace-ys/vault-init/workflows/release/badge.svg
[release-workflow]: https://github.com/jace-ys/vault-init/actions?query=workflow%3Arelease

# `vault-init`

`vault-init` is a small utility for automating the initialization and unsealing of [HashiCorp Vault](https://www.vaultproject.io/). It draws inspiration from [kelseyhightower/vault-init](https://github.com/kelseyhightower/vault-init), but doesn't rely on any public cloud infrastructure for the encrypting and storing of Vault's root tokens and unseal keys.

You would typically use this if you do not have access to public cloud infrastructure, or if your Vault deployment must operate entirely on-prem. For most production deployments, you would want to use Vault's [native auto-unsealing capabilities](https://www.vaultproject.io/docs/concepts/seal#auto-unseal) if possible.

## Overview

`vault-init` is written in Go and packaged as a binary that exposes a command-line interface. Its core is the [`start`](docs/examples.md#start) command that launches a daemon process designed to be run alongside a Vault server and communicate with it over localhost. It will continuously poll the status of the Vault server and depending on its state, automatically initialize and/or unseal it.

After `vault-init` initializes a Vault server, it encrypts the initial root token and unseal keys before storing them for future use in unsealing operations; this runs on the idea of pluggable backends for both encryption and storage that you can mix-and-match, [configurable through the CLI](#configuration).

For the full list of encryption and storage backends currently supported, see [`docs/encryption.md`](docs/encryption.md) and [`docs/storage.md`](docs/storage.md) respectively.

## Installation

#### Binary

Pre-compiled `vault-init` binaries for various platforms can be found under the [Releases](https://github.com/jace-ys/vault-init/releases) section of this repository.

#### Source

Clone this repository and build the binary from source using the given Makefile (requires `go 1.16+`):

```shell
$ make
```

This will compile and place the `vault-init` binary into a local `dist` directory.

#### Docker

A Docker image for `vault-init` is available on [Docker Hub](https://hub.docker.com/repository/docker/jaceys/vault-init) and can be pulled via:

```shell
$ docker pull docker.io/jaceys/vault-init:latest
```

## Usage

To use the `vault-init` CLI:

```shell
$ vault-init [<flags>] <command> [<args> ...]
```

#### Configuration

To view all configuration options of each command, use the `--help` flag:

```shell
$ vault-init --help
```

Configuration options can also be passed in as environment variables, using the uppercased snake-case version of the respective flag name (eg. `VAULT_ADDR` for `--vault-addr`).

Most commands require you to specify the encryption and storage backend to use via the `--encryption` and `--storage` flag, respectively. Each backend has its own set of configuration options, with their names typically following the given patterns:

```
--encryption-[backend-name]-[flag-name]
--storage-[backend-name]-[flag-name]
```

You will need to specify the appropriate flags depending on the backends you have chosen. Full documentation on configuration options for each backend can be found in [`docs/encryption.md`](docs/encryption.md) and [`docs/storage.md`](docs/storage.md).

## Examples

Examples on using `vault-init` can be found in [`docs/examples.md`](docs/examples.md).

## Contributing

All contributions are welcome, so if you don't see an encryption/storage backend that you would like to use, simply open an issue or pull request to propose it. Have a look at the code in [`pkg/encryption`](pkg/encryption) and [`pkg/storage`](pkg/storage) for ideas on how to contribute.

## License

See [LICENSE](LICENSE).
