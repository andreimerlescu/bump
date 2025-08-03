#!/bin/bash

[ "${IN_TEST_FILE}" != "1" ] && { echo "Cannot call this script outside of test.sh!"; exit 1; }

declare counterName="${counterName:-"bump-test-passes"}"

# Start with a populated VERSION file and bump its alpha
declare -a scenario_01=(
  "echo \"v1.0.0\" > VERSION"
  "bump -check"
  "cat VERSION"
  "grep 'v1.0.0' VERSION"
  "bump -alpha"
  "grep 'v1.0.0' VERSION"
  "cat VERSION"
  "bump -alpha -write"
  "cat VERSION"
  "grep 'v1.0.0-alpha.1' VERSION"
  "rm VERSION"
)

# Start with an alpha version 0 and work from there
declare -a scenario_02=(
  "echo \"v1.0.0-alpha.0\" > VERSION"
  "bump -check"
  "cat VERSION"
  "grep 'v1.0.0-alpha.0' VERSION"
  "bump -alpha"
  "cat VERSION"
  "grep 'v1.0.0-alpha.0' VERSION"
  "bump -alpha -write"
  "cat VERSION"
  "grep 'v1.0.0-alpha.1' VERSION"
  "bump -patch"
  "cat VERSION"
  "grep 'v1.0.0-alpha.1' VERSION"
  "bump -patch -write"
  "cat VERSION"
  "grep 'v1.0.1' VERSION"
  "bump -major -write"
  "cat VERSION"
  "grep 'v2.0.0' VERSION"
  "bump -preview -write"
  "grep 'v2.0.0-preview.1' VERSION"
  "rm VERSION"
)

# work with an invalid VERSION format and fix it
declare -a scenario_03=(
  "echo \"1.25\" > VERSION"
  "bump -fix"
  "grep '1.25' VERSION"
  "bump -fix -write"
  "grep 'v1.25.0' VERSION"
  "rm VERSION"
)

# Using bump in a manner that does not change the VERSION file at all
declare -a scenario_04=(
  "echo \"v1.17.7-beta.6\" > VERSION"
  "bump -check -fix"
  "cat VERSION"
  "grep 'v1.17.7-beta.6' VERSION"
  "bump -check -fix -write"
  "cat VERSION"
  "grep 'v1.17.7-beta.6' VERSION"
  "rm VERSION"
)

# go application version bump using igo version to correct
declare -a scenario_05=(
  "echo \"module testApp-${counterName}\" > go.mod"
  "echo \"\" >> go.mod"
  "echo \"go 1.24\" >> go.mod"
  "cat go.mod"
  "grep 'go 1.24' go.mod"
  "bump -in go.mod -fix"
  "grep 'go 1.24' go.mod"
  "bump -in go.mod -fix -write"
  "grep 'go $(cat ~/go/version)' go.mod"
  "cat go.mod"
  "rm go.mod"
)

declare -a scenario_06=(
  "echo v1.0.0 > VERSION"
  "bump -json -check"
  "cat VERSION"
  "bump -json -beta"
  "cat VERSION"
  "bump -json -beta -write"
  "cat VERSION"
  "rm VERSION"
)

declare -a scenario_07=(
  "echo '${empty_package_json}' | base64 -d | tee package.json > /dev/null"
  "cat package.json"
  "bump -in package.json -fix"
  "bump -in package.json -fix -write"
  "bump -in package.json -patch"
  "bump -in package.json -patch -write"
  "grep '\"version\": \"1.2.4\"' package.json"
  "bump -in package.json -json -minor"
  "bump -in package.json -minor -write"
  "grep '\"version\": \"1.3.0\"' package.json"
  "rm package.json"
)

declare -a scenario_08=(
  "echo 'module myapp' > go.mod && echo 'go 1.24' >> go.mod"
  "bump -in go.mod -check"
  "grep 'go 1.24' go.mod"
  "bump -in go.mod -patch"
  "grep 'go 1.24' go.mod"
  "bump -in go.mod -patch -write"
  "grep 'go $(cat ~/go/version)' go.mod"
  "bump -in go.mod -patch -write"
  "grep 'go $(cat ~/go/version)' go.mod"
  "bump -in go.mod -patch -write"
  "grep 'go $(cat ~/go/version)' go.mod"
  "bump -in go.mod -patch -write"
  "grep 'go $(cat ~/go/version)' go.mod"
  "bump -in go.mod -patch -write"
  "grep 'go $(cat ~/go/version)' go.mod"
  "rm go.mod"
)

declare -a scenario_09=(
  "echo 'LABEL version=\"v3.2.1\"' > Dockerfile"
  "bump -in Dockerfile -check"
  "grep 'LABEL version=\"v3.2.1\"' Dockerfile"
  "bump -in Dockerfile -patch"
  "grep 'LABEL version=\"v3.2.1\"' Dockerfile"
  "bump -in Dockerfile -patch -write"
  "grep 'LABEL version=\"v3.2.2\"' Dockerfile"
  "rm Dockerfile"
)

declare -a scenario_10=(
  "echo '${empty_chart_yaml}' | base64 -d | tee Chart.yaml > /dev/null"
  "bump -in Chart.yaml -check"
  "bump -in Chart.yaml -patch"
  "bump -in Chart.yaml -patch -write"
  "grep 'version: 0.1.1' Chart.yaml"
  "rm Chart.yaml"
)

declare -a scenario_11=(
  "echo '<project><version>2.2.2</version></project>' > pom.xml"
  "bump -in pom.xml -check"
  "bump -in pom.xml -patch"
  "bump -in pom.xml -patch -write"
  "grep '<version>2.2.3</version>' pom.xml"
  "rm pom.xml"
)

declare -a scenario_12=(
  "echo \"v5.5.5\" > VERSION"
  "BUMP_ALWAYS_WRITE=true bump -env"
  "BUMP_ALWAYS_WRITE=true bump -patch"
  "grep 'v5.5.6' VERSION"
  "BUMP_DEFAULT_INPUT=VERSION bump -minor"
  "BUMP_DEFAULT_INPUT=VERSION bump -minor -write"
  "grep 'v5.6.0' VERSION"
  "rm VERSION"
)

declare -a scenario_13=(
  "bump -parse v1.2.3-alpha.4 -init"
  "cat VERSION"
  "bump -parse v2.3.4-alpha.5 -write"
  "cat VERSION"
  "rm VERSION"
  "bump -parse v3.4.5-alpha.6 -init"
  "cat VERSION"
  "rm VERSION"
)

declare -a scenario_14=(
  "echo '${populated_package_json}' | base64 -d | tee package.json > /dev/null"
  "cat package.json"
  "bump -in package.json -check"
  "bump -in package.json -fix"
  "bump -in package.json -fix -write"
  "bump -in package.json -patch"
  "bump -in package.json -patch -write"
  "grep '\"version\": \"1.0.1\"' package.json"
  "bump -in package.json -json -minor"
  "bump -in package.json -minor -write"
  "grep '\"version\": \"1.1.0\"' package.json"
  "cat package.json"
  "rm package.json"
)

export scenario_01
export scenario_02
export scenario_03
export scenario_04
export scenario_05
export scenario_06
export scenario_07
export scenario_08
export scenario_09
export scenario_10
export scenario_11
export scenario_12
export scenario_13
export scenario_14
