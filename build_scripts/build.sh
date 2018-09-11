#!/usr/bin/env bash

set -e
set -u

BINARIES_FOLDER="binaries"

echo "
###
## Build binaries
###"
gox -os="linux darwin" -arch="amd64" -output="${BINARIES_FOLDER}/vault-iam-request_{{.OS}}-{{.Arch}}"


echo "
###
## UPX binaries
###"
upx "${BINARIES_FOLDER}"/*



echo "
###
## Tar everything up
###"

cd "${BINARIES_FOLDER}"
for binary in vault-*; do
    IFS="_" read -a tokens <<< "${binary}"
    echo "Compressing ${binary}"
    mv "${binary}" "${tokens[0]}"
    tar -czf "${tokens[0]}-${CIRCLE_TAG}-${tokens[1]}.tar.gz" "${tokens[0]}"
    rm "${tokens[0]}"
done
cd ..

ls "${BINARIES_FOLDER}"

echo "
###
## Upload release to Github
###
"
ghr -n "${CIRCLE_TAG}" "${CIRCLE_TAG}" "${BINARIES_FOLDER}/"