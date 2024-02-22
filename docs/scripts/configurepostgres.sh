#!/usr/bin/env bash
set -e

LAO = "#listen_addresses = 'localhost'"
LAR = "listen_addresses = '*'"
apt-get install sudo
  
sudo -u postgres psql -c "CREATE USER keycloak WITH PASSWORD 'keycloak';"
sudo -u postgres psql -c "ALTER USER postgres WITH PASSWORD 'postgres';"
sudo -u postgres createdb -O keycloak -E 'UTF-8' keycloak
sudo -u postgres sed -i 's/$LAO/$LAR/g' /etc/postgresql/16/main/postgresql.conf
sudo -u postgres echo "host    all             all             192.168.56.1/32         scram-sha-256" >>/etc/postgresql/16/main/pg_hba.conf
sudo -u postgres echo "host    all             all             127.0.0.1/32            scram-sha-256" >>/etc/postgresql/16/main/pg_hba.conf
