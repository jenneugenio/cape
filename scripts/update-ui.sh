#!/bin/sh

if [ -z "${CAPE_UI_BUCKET}" ]
then
	echo 'Must specify location of UI assets in CAPE_UI_BUCKET'
	exit 1
fi

GIT_ROOT="$(git rev-parse --show-toplevel)"

gsutil cat "gs://${CAPE_UI_BUCKET}/manifest.json" | jq -r .latest | tee "${GIT_ROOT}/UI_VERSION"

