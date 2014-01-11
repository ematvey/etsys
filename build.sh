#!/bin/bash
export GOPATH="`pwd`"
export PATH="$PATH:`pwd`/bin"

#go get github.com/lib/pq

go install etsysd
go install simulator

#go install sqlc
