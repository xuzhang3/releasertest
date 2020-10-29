#!/usr/bin/env bash

set -euo pipefail

SCRIPTS_DIR="$(cd -P "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BUILD_DIR="${SCRIPTS_DIR}/../dist/"
SOURCE_DIR="${SCRIPTS_DIR}/../"
NAME="azuredevops"
VERSION=$(git tag | sort -V | tail -1)

#set API_JSON=$(printf '{"tag_name": "v%s","target_commitish": "master","name": "v%s","body": "Release of version %s","draft": false,"prerelease": false}' $1 $1 $1)
#curl --data "$API_JSON" https://api.github.com/repos/:owner/:repository/releases?access_token=:access_token

OS_ARCH=( "cow:moo"
        "dinosaur:roar"
        "bird:chirp"
        "bash:rock" )


function clean() {
  info "Cleaning $BUILD_DIR"
  rm -rf "$BUILD_DIR"
  mkdir -p "$BUILD_DIR"
}

function release() {
  echo "Clean build directory"
  clean
  echo $(zip --help)
  for os_arch in "${OS_ARCH[@]}" ; do
    KEY=${OS_ARCH%%:*}
    VALUE=${OS_ARCH#*:}
    info "%s likes to %s.\n" "$KEY" "$VALUE"
  done
#
#  BUILD_ARTIFACT="terraform-provider-${NAME}_v${VERSION}"
#  BUILD_ARTIFACT_ZIP="$BUILD_ARTIFACT_${VERSION}_${OS}_${ARCH}.zip"
#  info "Attempting to build $BUILD_ARTIFACT"
#  (
#    cd "$SOURCE_DIR"
#    go mod download
#    go build -o "$BUILD_DIR/$BUILD_ARTIFACT"
#    zip $BUILD_ARTIFACT_ZIP $BUILD_ARTIFACT
#  )

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


