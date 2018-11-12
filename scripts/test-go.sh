#!/usr/bin/env bash

set -euo pipefail

echo
echo "==> Running automated tests <=="
go test -coverprofile=coverage.out ./apiclient ./column ./record ./response ./responseparser ./responsetemplate ./responsetemplatemanager ./socketconfig
go tool cover -html=coverage.out -o coverage.html
cd .. || exit
exit
