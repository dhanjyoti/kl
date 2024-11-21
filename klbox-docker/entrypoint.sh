#! /usr/bin/env bash

set -o errexit
set -o pipefail

# Ensure SSH keys are generated if they don't exist
if [ ! -f /etc/ssh/keys/ssh_host_rsa_key ]; then
    ssh-keygen -t rsa -b 4096 -f /etc/ssh/keys/ssh_host_rsa_key -N ""
    ssh-keygen -t rsa -b 4096 -f /etc/ssh/keys/ssh_host_rsa_key.pub -N ""
fi
if [ ! -f /etc/ssh/keys/ssh_host_ecdsa_key ]; then
    ssh-keygen -t ecdsa -b 521 -f /etc/ssh/keys/ssh_host_ecdsa_key -N ""
    ssh-keygen -t ecdsa -b 521 -f /etc/ssh/keys/ssh_host_ecdsa_key.pub -N ""
fi
if [ ! -f /etc/ssh/keys/ssh_host_ed25519_key ]; then
    ssh-keygen -t ed25519 -f /etc/ssh/keys/ssh_host_ed25519_key -N ""
    ssh-keygen -t ed25519 -f /etc/ssh/keys/ssh_host_ed25519_key.pub -N ""
fi

chmod 600 /etc/ssh/keys/ssh_host_*
/usr/sbin/sshd -D -p "$SSH_PORT" &
pid=$!

cat >/kl-tmp/kill-sshd.sh <<EOF
sudo kill $@ $pid
EOF

wait $pid
