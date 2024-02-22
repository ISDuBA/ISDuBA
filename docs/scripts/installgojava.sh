#!/usr/bin/env bash
set -e

apt-get update

apt install -y openjdk-17-jre-headless

wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz

rm  /usr/local/go -rf && tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz

rm go1.22.0.linux-amd64.tar.gz

ln -snf /usr/local/go/bin/go /usr/local/bin/go
