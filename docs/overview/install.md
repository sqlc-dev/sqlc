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

Get pre-built binaries for *v1.16.0*:

- [Linux](https://github.com/kyleconroy/sqlc/releases/download/v1.16.0/sqlc_1.16.0_linux_amd64.tar.gz)
- [macOS](https://github.com/kyleconroy/sqlc/releases/download/v1.16.0/sqlc_1.16.0_darwin_amd64.zip)
- [Windows (MySQL only)](https://github.com/kyleconroy/sqlc/releases/download/v1.16.0/sqlc_1.16.0_windows_amd64.zip)

See [downloads.sqlc.dev](https://downloads.sqlc.dev/) for older versions.
