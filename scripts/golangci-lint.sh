#!/usr/bin/env bash
set -e

if [ $LOCAL_RUN == "true" ]; then
	echo "Skipping golangci-lint installation. Local run enabled."
	exit 0
fi

# shellcheck disable=SC1091 # Not following.
. "$(dirname "$0")"/common.sh

if which golangci-lint; then
	echo "golint installed"
else
	echo "Downloading golint tool"

	if [[ -z "${GOPATH}" ]]; then
    DEPLOY_PATH=/tmp/
  else
    DEPLOY_PATH=${GOPATH}/bin
  fi
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "${DEPLOY_PATH}" v1.51.1
fi

PATH=${PATH}:${DEPLOY_PATH} golangci-lint run -v
