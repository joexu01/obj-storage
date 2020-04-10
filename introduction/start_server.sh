#!/bin/bash
export STORAGE_ROOT=$HOME/tmp
export LISTEN_ADDRESS=:8000
go run ./server.go
