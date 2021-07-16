#!/bin/bash
# Usage:
#  cd cosmos-sdk
#  ~/certus/janus/patches/cosmos-sdk/dump_patches.sh upstream/v0.29.1

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"

rm ${DIR}/*.patch

git format-patch $1 -o ${DIR}

echo "$(git describe --tags)" > ${DIR}/ref
echo "$(git describe --tags $1)" > ${DIR}/upstream_ref

echo "Wrote patch files to ${DIR}."
