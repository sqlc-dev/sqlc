# Uploading projects

*Added in v1.22.0*

Uploading an archive of your project ensures that future releases of sqlc do not
break your code. Similar to Rust's [crater](https://github.com/rust-lang/crater)
project, uploaded archives are tested against development releases of sqlc to
verify correctness.

Interested in uploading projects? Sign up [here](https://docs.google.com/forms/d/e/1FAIpQLSdxoMzJ7rKkBpuez-KyBcPNyckYV-5iMR--FRB7WnhvAmEvKg/viewform) or send us an email
at [hello@sqlc.dev](mailto:hello@sqlc.dev).

## Add configuration

After creating a project, add the project ID to your sqlc configuration file.

```yaml
version: "2"
cloud:
  project: "<PROJECT-ID>"
```

You'll also need to create an auth token and make it available via the
`SQLC_AUTH_TOKEN` environment variable.

```shell
export SQLC_AUTH_TOKEN=sqlc_xxxxxxxx
```

## Dry run

You can see what's included when uploading your project by using using the
`--dry-run` flag:

```shell
sqlc upload --dry-run
```

The output is the request `sqlc` would have sent without the `--dry-run` flag.

## Upload

Once you're ready to upload, remove the `--dry-run` flag.

```shell
sqlc upload
```

By uploading your project, you're making sqlc more stable and reliable. Thanks!