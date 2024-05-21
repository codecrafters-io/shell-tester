#!/bin/sh
# echo "hey"
set -e
cd $(dirname "$0")
tmpFile=$(mktemp)
go build -o "$tmpFile" main.go
exec $tmpFile
