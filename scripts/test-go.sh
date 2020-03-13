#!/usr/bin/env bash
echo
echo "==> Running automated tests <=="
go test -v -race -coverprofile=coverage.out ./apiclient ./column ./record ./response ./responseparser ./responsetemplate ./responsetemplatemanager ./socketconfig
go tool cover -html=coverage.out -o coverage.html
