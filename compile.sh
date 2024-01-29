#! /usr/bin/env nix-shell
#! nix-shell -i bash -p go bash

CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o ./job_share
