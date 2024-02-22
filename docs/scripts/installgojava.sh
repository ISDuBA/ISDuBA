#!/usr/bin/env bash
set -e

apt install -y openjdk-17-jre-headless unzip

wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz

rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz

rm -rf go1.22.0.linux-amd64.tar.gz

echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
