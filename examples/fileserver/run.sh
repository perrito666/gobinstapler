#!/bin/sh
set -e
go build -o thinselfserve .
gobinstapler thinselfserve selfserve $@
