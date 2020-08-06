#!/usr/bin/env bash
set -x

for config in internal/endtoend/testdata/**/sqlc.json; do (cd $(dirname "$config") && sqlc-dev generate); done
for dir in examples/*; do (cd "$dir" && sqlc-dev generate); done
