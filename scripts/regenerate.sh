#!/usr/bin/env bash
set -x

for dir in internal/endtoend/testdata/*; do (cd "$dir" && sqlc-dev generate); done
for dir in examples/*; do (cd "$dir" && sqlc-dev generate); done
