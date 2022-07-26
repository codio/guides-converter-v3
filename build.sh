#!/bin/bash
set -xe

commit=$1
out_dir=$2

output_name="guides-converter-v3-$commit"
go build -ldflags "-s -w" -o $out_dir/$output_name ./cmd/guides-converter-v3
