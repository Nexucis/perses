#!/bin/bash

set -e

cd ui/

files=("../LICENSE" "../CHANGELOG.md")
workspaces=$(cat package.json | jq -r '.workspaces[]')

function copy() {
  for file in "${files[@]}"; do
    for workspace in ${workspaces}; do
      if [ -f "${file}" ]; then
        cp "${file}" "${workspace}"/"$(basename "${file}")"
      fi
    done
  done
}

function publish() {
  dry_run="${1}"
  cmd="npm publish --access public"
  if [[ "${dry_run}" == "dry-run" ]]; then
    cmd+=" --dry-run"
  fi
  for workspace in ${workspaces}; do
    # package "app" is private so we shouldn't try to publish it.
    if [[ "${workspace}" != "app" ]]; then
      cd "${workspace}"
      eval "${cmd}"
      cd ../
    fi
  done

}

function checkVersion() {
  version=${1}
  if [[ "${version}" =~ ^v[0-9]+(\.[0-9]+){2}(-.+)?$ ]]; then
    echo "version '${version}' follows the semver"
  else
    echo "version '${version}' doesn't follow the semver"
    exit 1
  fi
}

function checkPackage() {
  version=${1}
  if [[ "${version}" == v* ]]; then
    version="${version:1}"
  fi
  for workspace in ${workspaces}; do
    cd "${workspace}"
    package_version=$(npm run env | grep npm_package_version | cut -d= -f2-)
    if [ "${version}" != "${package_version}" ]; then
      echo "version of ${workspace} is not the correct one"
      echo "expected one: ${version}"
      echo "current one: ${package_version}"
      echo "please use ./ui_release --release ${version}"
      exit 1
    fi
    cd ..
  done
}

function clean() {
  for file in "${files[@]}"; do
    for workspace in ${workspaces}; do
      f="${workspace}"/"$(basename "${file}")"
      if [ -f "${f}" ]; then
        rm "${f}"
      fi
    done
  done
}

function bumpVersion() {
  version="${1}"
  if [[ "${version}" == v* ]]; then
    version="${version:1}"
  fi
  # increase the version on all packages
  npm version "${version}" --workspaces
  # upgrade the @perses-dev/* dependencies on all packages
  for workspace in ${workspaces}; do
    sed -E -i "" "s|(\"@perses-dev/.+\": )\".+\"|\1\"\^${version}\"|" "${workspace}"/package.json
  done
}

function tag() {
  version="${1}"
  tag="v${version}"
  branch=$(git branch --show-current)
  checkVersion "${tag}"
  expectedBranch="release/$(echo "${tag}" | sed -E  's/(v[0-9]+\.[0-9]+).*/\1/')"

  if [[ "${branch}" != "${expectedBranch}" ]]; then
    echo "you are not on the correct release branch (i.e. not on ${expectedBranch}) to create the tag"
    exit 1
  fi

  git pull origin "${expectedBranch}"
  git tag -s "${tag}" -m "${tag}"
}

if [[ "$1" == "--copy" ]]; then
  copy
fi

if [[ $1 == "--publish" ]]; then
  publish "${@:2}"
fi

if [[ $1 == "--check-package" ]]; then
  checkPackage "${@:2}"
fi

if [[ $1 == "--check-version" ]]; then
  checkVersion "${@:2}"
fi

if [[ $1 == "--bump-version" ]]; then
  bumpVersion "${@:2}"
fi

if [[ $1 == "--clean" ]]; then
  clean
fi

if [[ $1 == "--tag" ]]; then
  tag "${@:2}"
fi