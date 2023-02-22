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

## go install 

### Go >= 1.17:

```
go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
```

### Go < 1.17:

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

Run `sqlc` using `docker run` in the Command Prompt on Windows (`cmd`):

```
docker run --rm -v "%cd%:/src" -w /src kjconroy/sqlc generate
```

## Downloads

Get pre-built binaries for *v1.17.2*:

- [Linux](https://github.com/kyleconroy/sqlc/releases/download/v1.17.2/sqlc_1.17.2_linux_amd64.tar.gz)
- [macOS](https://github.com/kyleconroy/sqlc/releases/download/v1.17.2/sqlc_1.17.2_darwin_amd64.zip)
- [Windows (MySQL only)](https://github.com/kyleconroy/sqlc/releases/download/v1.17.2/sqlc_1.17.2_windows_amd64.zip)

See [downloads.sqlc.dev](https://downloads.sqlc.dev/) for older versions.
