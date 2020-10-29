#!/usr/bin/env bash

set -euo pipefail

SCRIPTS_DIR="$(cd -P "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BUILD_DIR="${SCRIPTS_DIR}/dist/"
SOURCE_DIR="${SCRIPTS_DIR}/"
NAME="azuredevops"
BUILD_ARTIFACT="terraform-provider-${NAME}_${VERSION}"


#set API_JSON=$(printf '{"tag_name": "v%s","target_commitish": "master","name": "v%s","body": "Release of version %s","draft": false,"prerelease": false}' $1 $1 $1)
#curl --data "$API_JSON" https://api.github.com/repos/:owner/:repository/releases?access_token=:access_token

OS_ARCH=("freebsd:amd64"
  "freebsd:386"
  "freebsd:arm"
  "freebsd:arm64"
  "windows:amd64"
  "windows:386"
  "linux:amd64"
  "linux:386"
  "linux:arm"
  "linux:arm64"
  "darwin:amd64")


function clean() {
  info "Cleaning $BUILD_DIR"
  rm -rf "$BUILD_DIR"
  mkdir -p "$BUILD_DIR"
}

function release() {
  info "Clean build directory"
  clean

  info "Attempting to build ${BUILD_ARTIFACT}"

  cd "$SOURCE_DIR"
  go mod download
  for os_arch in "${OS_ARCH[@]}" ; do
    OS=${os_arch%%:*}
    ARCH=${os_arch#*:}
    info "GOOS: ${OS}, GOARCH: ${ARCH}"
    (
      env GOOS="${OS}" GOARCH="${ARCH}" go build -o "${BUILD_DIR}${BUILD_ARTIFACT}"
      zip -r "${BUILD_DIR}${BUILD_ARTIFACT}_${OS}_${ARCH}.zip" "${BUILD_DIR}${BUILD_ARTIFACT}"
#      tar -cf "${BUILD_DIR}${BUILD_ARTIFACT}_${OS}_${ARCH}.tar" "${BUILD_DIR}${BUILD_ARTIFACT}"
    )
  done
  ls -al ${BUILD_DIR}*.zip
  shasum -a 256 ${BUILD_DIR}*.zip > "./dist/${BUILD_ARTIFACT}_SHA256SUMS"
  ls -al ${BUILD_DIR}*.zip

}

function log() {
    LEVEL="$1"
    shift
    echo "[$LEVEL] $@"
}

function info() {
    log "INFO" $@
}

function fatal() {
    log "FATAL" $@
    exit 1
}

release
