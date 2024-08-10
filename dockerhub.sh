#!/usr/bin/env bash

set -o nounset -o pipefail -o errexit

if [[ ${TRACE:-0} == "1" ]]; then
  set -o xtrace
fi

script_dir() {
  local FILE_SOURCE="${BASH_SOURCE[0]}"

  if [[ -L ${FILE_SOURCE} ]]; then
    dirname "$(readlink "${FILE_SOURCE}")"
  else
    (
      cd "$(dirname "${FILE_SOURCE}")" && pwd
    )
  fi
}

main() {
  local SCRIPT_DIR
  SCRIPT_DIR="$(script_dir)"

  source "${SCRIPT_DIR}/scripts/meta" && meta_check "var" "pass"

  var_read DOCKER_OWNER "${1-}"
  var_read DOCKER_USER "$(pass_get "dev/docker" "login")"
  var_read DOCKER_PASSWORD "$(pass_get "dev/docker" "password")" "secret"

  for repo in $(go run "${SCRIPT_DIR}" -username "${DOCKER_USER}" -password "${DOCKER_PASSWORD}" -owner "${DOCKER_OWNER}" -list); do
    go run "${SCRIPT_DIR}" -username "${DOCKER_USER}" -password "${DOCKER_PASSWORD}" -image "${repo}" -delete -grep '^[a-f0-9]{7,8}($|-)'
    go run "${SCRIPT_DIR}" -username "${DOCKER_USER}" -password "${DOCKER_PASSWORD}" -image "${repo}" -delete -grep '[0-9]{12}' -last
    go run "${SCRIPT_DIR}" -username "${DOCKER_USER}" -password "${DOCKER_PASSWORD}" -image "${repo}" -delete -grep '^v(?P<tagBucket>\d+.\d+)(?:.\d+)' -last
  done

  unset DOCKER_PASSWORD
}

main "${@}"
