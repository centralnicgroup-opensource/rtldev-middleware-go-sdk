#!/usr/bin/env zsh
# shellcheck shell=bash
# Runs after the devbase Feature's own post-create. Everything the previous
# version of this script did besides this -- installing commitizen globally,
# symlinking ~/.zsh_history out of /WSL_USER, writing ~/.czrc -- is devbase's
# job or a host mount now.
set -e

echo "=> Running go mod tidy"
go mod tidy
