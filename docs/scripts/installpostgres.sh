#!/usr/bin/env bash
set -e

apt install vim gnupg2 -y
curl -fsSL https://www.postgresql.org/media/keys/ACCC4CF8.asc| gpg --trust-model always --dearmor -o /etc/apt/trusted.gpg.d/postgresql.gpg
sh -c 'echo "deb https://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list' 
apt update
apt install postgresql-16 -y
