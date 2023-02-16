#!/usr/bin/env bash

if [ $LOCAL_RUN == "true" ]; then
	echo "Skipping ginkgo installation. Local run enabled."
	exit 0
fi

# Only overwrite the GOPATH if needed.
if [ -z $GOPATH ]; then
	GOPATH="/root/go"
	export PATH=$PATH:$GOPATH/bin
fi

GINKGO_OLD_VERSION="Ginkgo Version 1.16.5"

if ! which ginkgo || ginkgo version -eq "$GINKGO_OLD_VERSION"; then {  
	echo "Downloading ginkgo tool"
	go install "$(awk '/ginkgo/ {printf "%s/ginkgo@%s", $1, $2}' go.mod)"
} fi
