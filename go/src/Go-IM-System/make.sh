#!/bin/bash
rm -rf async
rm -rf client 
go build -o async main.go server.go user.go
go build -o client client.go
