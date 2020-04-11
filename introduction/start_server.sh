#!/bin/bash

if [ ! -d "$HOME/tmp/objects" ]; then
  echo "cannot specify storage path; creating a folder at $HOME/tmp/objects automatically"
  mkdir "$HOME/tmp/objects"
fi

echo "storage path checked... starting the server; listening at :8000"

export STORAGE_ROOT=$HOME/tmp
export LISTEN_ADDRESS=:8000
go run ./server.go
