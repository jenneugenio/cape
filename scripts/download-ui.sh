#!/bin/sh

if [ -z "${CAPE_UI_BUCKET}" ]
then
	echo 'Must specify location of UI assets in CAPE_UI_BUCKET'
	exit 1
fi

GIT_ROOT="$(git rev-parse --show-toplevel)"
CAPE_UI_VERSION=${CAPE_UI_VERSION:-"$(tr -d '\n' < "${GIT_ROOT}/UI_VERSION")"}
CAPE_UI_DIR=${CAPE_UI_DIR:-"coordinator/ui/assets"}

mkdir -p "${CAPE_UI_DIR}"

echo "Downloading UI version ${CAPE_UI_VERSION}..."

if ( gsutil cat "gs://${CAPE_UI_BUCKET}/${CAPE_UI_VERSION}.tar" | tar -zx -C "${CAPE_UI_DIR}" -f - )
then
	echo "done"
else
	echo "failed"
fi
