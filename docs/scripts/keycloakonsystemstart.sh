#!/usr/bin/env bash
set -e

echo "[Unit]
Description=Keycloak
After=network.target

[Service]
Type=idle
User=keycloak
Group=keycloak
ExecStart=/opt/keycloak/bin/kc.sh start-dev
TimeoutStartSec=600
TimeoutStopSec=600

[Install]
WantedBy=multi-user.target
" > /etc/systemd/system/keycloak.service

systemctl enable keycloak
systemctl start keycloak
