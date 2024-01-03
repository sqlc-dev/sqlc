# `push` - Uploading projects

```{note}
`push` is powered by [sqlc Cloud](https://dashboard.sqlc.dev). Sign up for [free](https://dashboard.sqlc.dev) today.
```

*Added in v1.24.0*

We've renamed the `upload` sub-command to `push`. We've also changed the data sent along in a push request. Upload used to include the configuration file, migrations, queries, and all generated code. Push drops the generated code in favor of including the [plugin.GenerateRequest](https://buf.build/sqlc/sqlc/docs/main:plugin#plugin.GenerateRequest), which is the protocol buffer message we pass to codegen plugins.

## Add configuration

After creating a project, add the project ID to your sqlc configuration file.

```yaml
version: "2"
cloud:
  project: "<PROJECT_ID>"
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
$ sqlc push --dry-run
2023/11/21 10:39:51 INFO config file=sqlc.yaml bytes=912
2023/11/21 10:39:51 INFO codegen_request queryset=app file=codegen_request.pb
2023/11/21 10:39:51 INFO schema queryset=app file=migrations/00001_initial.sql bytes=3033
2023/11/21 10:39:51 INFO query queryset=app file=queries/app.sql bytes=1150
```

The output is the files `sqlc` would have sent without the `--dry-run` flag.

## Push

Once you're ready to push, remove the `--dry-run` flag.

```shell
$ sqlc push
```

### Tags

You can provide tags to associate with a push, primarily as a convenient reference when using `sqlc verify` with the `against` argument.

Tags only refer to a single push, so if you pass an existing tag to `push` it will overwrite the previous reference.

```shell
$ sqlc push --tag main
```

### Annotations 

Annotations are added to each push request. By default, we include these environment variables (if they are present).

```
GITHUB_REPOSITORY
GITHUB_REF
GITHUB_REF_NAME
GITHUB_REF_TYPE
GITHUB_SHA
```
