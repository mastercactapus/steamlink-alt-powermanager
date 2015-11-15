#!/usr/bin/env bash
set -e

GOARCH=arm GOARM=7 go build -o powermanager
