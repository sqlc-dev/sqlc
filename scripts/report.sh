#!/usr/bin/env bash

curl \
  -X POST \
  --fail-with-body \
  -H "Authorization: Token token=\"$BUILDKITE_ANALYTICS_TOKEN\"" \
  -F "data=@junit.xml" \
  -F "format=junit" \
  -F "run_env[CI]=buildkite" \
  -F "run_env[key]=$GITHUB_RUN_ID" \
  -F "run_env[number]=$GITHUB_RUN_NUMBER-$GITHUB_RUN_ATTEMPT" \
  -F "run_env[job_id]=$GITHUB_RUN_ID" \
  -F "run_env[branch]=$GITHUB_REF" \
  -F "run_env[commit_sha]=$GITHUB_SHA" \
  -F "run_env[message]=Foo" \
  -F "run_env[url]=$GITHUB_SERVER_URL" \
  https://analytics-api.buildkite.com/v1/uploads
