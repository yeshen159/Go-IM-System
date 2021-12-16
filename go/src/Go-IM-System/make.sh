#!/bin/bash
rm -rf async
go build -o async main.go server.go user.go
