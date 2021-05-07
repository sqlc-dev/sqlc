# Installing sqlc

sqlc is distributed as a single binary with zero dependencies.

## macOS

```
brew install sqlc
```

## Ubuntu

```
sudo snap install sqlc
```

## go get

```
go get github.com/kyleconroy/sqlc/cmd/sqlc
```

## Docker

```
docker pull kjconroy/sqlc
```

Run `sqlc` using `docker run`:

```
docker run --rm -v $(pwd):/src -w /src kjconroy/sqlc generate
```

## Downloads

Get pre-built binaries for *v1.8.0*:

- [Linux](https://github.com/kyleconroy/sqlc/releases/download/v1.8.0/sqlc-v1.8.0-linux-amd64.tar.gz)
- [macOS](https://github.com/kyleconroy/sqlc/releases/download/v1.8.0/sqlc-v1.8.0-darwin-amd64.zip)
- [Windows (MySQL only)](https://github.com/kyleconroy/sqlc/releases/download/v1.8.0/sqlc-v1.8.0-windows-amd64.zip)

Binaries for a specific release can be downloaded on
[GitHub](https://github.com/kyleconroy/sqlc/releases).


## Tip Releases

Each commit is deployed to the [`devel` channel on Equinox](https://dl.equinox.io/sqlc/sqlc/devel):

- [Linux](https://bin.equinox.io/c/gvM95th6ps1/sqlc-devel-linux-amd64.tgz)
- [macOS](https://bin.equinox.io/c/gvM95th6ps1/sqlc-devel-darwin-amd64.zip)
- [Windows (MySQL only)](https://bin.equinox.io/c/gvM95th6ps1/sqlc-devel-windows-amd64.zip)
