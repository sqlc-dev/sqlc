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

Binaries for a given release can be downloaded from the [stable channel on
Equinox](https://dl.equinox.io/sqlc/sqlc/stable) or the latest [GitHub
release](https://github.com/kyleconroy/sqlc/releases).

## Tip Releases

Each commit is deployed to the [`devel` channel on Equinox](https://dl.equinox.io/sqlc/sqlc/devel):

- [Linux](https://bin.equinox.io/c/gvM95th6ps1/sqlc-devel-linux-amd64.tgz)
- [macOS](https://bin.equinox.io/c/gvM95th6ps1/sqlc-devel-darwin-amd64.zip)
- [Windows](https://bin.equinox.io/c/gvM95th6ps1/sqlc-devel-windows-amd64.zip)
