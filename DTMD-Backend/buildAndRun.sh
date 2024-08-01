#!/bin/bash
HOST=192.168.0.74
go build -o api main.go
chmod +x api
./api