set -e
set -x

export CAPE_NAME=cape-user
export CAPE_EMAIL=cape_user@mycape.com
export CAPE_PASSWORD=capecape

if ! cape setup local http://localhost:8080; then
    cape config clusters remove -y local
    cape setup local http://localhost:8080
fi
