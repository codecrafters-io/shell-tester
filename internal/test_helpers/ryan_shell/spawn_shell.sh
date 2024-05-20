#!/bin/sh
cd $(dirname $0)
tmpFile=$(mktemp)
go build -o $tmpFile main.go
echo "Welcome to Ryan's shell!"
exec $tmpFile
