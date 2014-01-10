#!/bin/bash
export GOPATH="`pwd`"
export PATH="$PATH:`pwd`/bin"

go install etsysd
go install simulator
