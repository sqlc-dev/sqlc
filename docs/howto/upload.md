# Uploading projects

*This feature requires signing up for [sqlc Cloud](https://app.sqlc.dev), which is currently in beta.*

Uploading your project ensures that future releases of sqlc do not break your
existing code. Similar to Rust's [crater](https://github.com/rust-lang/crater)
project, uploaded projects are tested against development releases of sqlc to
verify correctness.

## Add configuration

After creating a project, add the project ID to your sqlc configuration file.

```yaml
version: "1"
project:
  id: "<PROJECT-ID>"
packages: []
```

```json
{
  "version": "1",
  "project": {
    "id": "<PROJECT-ID>"
  },
  "packages": [
  ]
}
```

You'll also need to create an API token and make it available via the
`SQLC_AUTH_TOKEN` environment variable.

```shell
export SQLC_AUTH_TOKEN=sqlc_xxxxxxxx
```

## Dry run

You can see what's included when uploading your project by using using the `--dry-run` flag:

```shell
sqlc upload --dry-run
```

The output will be the exact HTTP request sent by `sqlc`.

## Upload

Once you're ready to upload, remove the `--dry-run` flag.

```shell
sqlc upload
```

By uploading your project, you're making sqlc more stable and reliable. Thanks!
