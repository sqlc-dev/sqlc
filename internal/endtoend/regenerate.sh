#!/usr/bin/env bash
set -x
for dir in testdata/*; do (cd "$dir" && sqlc-dev generate); done
