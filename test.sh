#!/bin/bash

function safe_exit() {
    local msg="${1}"
    echo -e "${msg}"
    exit 1
}

declare NO_COLOR
NO_COLOR="${NO_COLOR:-false}"

go install github.com/andreimerlescu/counter@latest || safe_exit "failed to install counter"

[ ! -d .counters ] && { mkdir -p .counters || safe_exit "failed to mkdir .counters"; }

COUNTER_DIR=$(realpath .counters)

export COUNTER_USE_FORCE=1
export COUNTER_ALWAYS_YES=1

declare COUNTER_DIR
export COUNTER_DIR

declare counterName
counterName="bump-test-passes"
counter -name $counterName -reset -yes 1> /dev/null || safe_exit "failed to reset counter"

echo "Preparing test env..."
mkdir -p test-data || safe_exit "failed to mkdir test-data"
go build -o ./test-data/bump . || safe_exit "failed to build bump"
cd test-data || safe_exit "failed to cd test-data"
chmod +x bump || safe_exit "failed to chmod bump"

function red() { if [ "${NO_COLOR}" == "false" ]; then echo -e "\033[0;31m${1}\033[0m"; else echo "${1}"; fi; }
function blue() { if [ "${NO_COLOR}" == "false" ]; then echo -e "\033[0;34m${1}\033[0m";  else echo "${1}"; fi;}
function green() { if [ "${NO_COLOR}" == "false" ]; then echo -e "\033[0;32m${1}\033[0m";  else echo "${1}"; fi;}
function purple() { if [ "${NO_COLOR}" == "false" ]; then echo -e "\033[0;35m${1}\033[0m";  else echo "${1}"; fi;}
function yellow() { if [ "${NO_COLOR}" == "false" ]; then echo -e "\033[1;33m${1}\033[0m";  else echo "${1}"; fi;}
function pink() { if [ "${NO_COLOR}" == "false" ]; then echo -e "\033[1;36m${1}\033[0m";  else echo "${1}"; fi;}

function run() {
  local cmd="${1}"
  local -i testNo
  testNo=$(counter -name $counterName -add || safe_exit "failed to increase counter")
  local -i exitCode
  local output

  # Capture both stdout and stderr together
  output=$(eval "${cmd}" 2>&1)
  exitCode=$?

  local prefix
  prefix="$(purple "$(whoami)")@$(yellow bump):$(blue "$(basename "${0}")")"

  if (( exitCode == 0 )); then
    printf "%s ⚡ %s ⇒  %s\n" "${prefix}" "$(pink "Test #${testNo}")" "$(green "${cmd}")"
    [[ -n "$output" ]] && printf "%s\n" "$output"
  else
    printf "%s ⚡ %s ⇒  %s\n" "${prefix}" "$(pink "Test #${testNo}")" "$(red "${cmd}")"
    safe_exit "Command failed with exit code ${exitCode}:\noutput: ${output}"
  fi
}

function new_package_json() {
  input="${1}"
  output="${2}"
  json="{"
  while IFS=',' read -ra pairs; do
      for pair in "${pairs[@]}"; do
          key="${pair%=*}"
          value="${pair#*=}"
          json+="\"$key\":\"$value\","
      done
  done <<< "$input"
  json="${json%,}}"
  echo -e "echo $json $output"
}

declare -a tests=(
 "echo \"--- SCENARIO ONE ---\""
 "echo \"v1.0.0\" > VERSION"
 "bump -check"
 "cat VERSION"
 "bump -alpha"
 "cat VERSION"
 "bump -alpha -write"
 "cat VERSION"
 "rm VERSION"
 "echo \"--- SCENARIO TWO ---\""
 "echo \"v1.0.0-alpha.0\" > VERSION"
 "bump -check"
 "cat VERSION"
 "bump -alpha"
 "cat VERSION"
 "bump -alpha -write"
 "cat VERSION"
 "bump -patch"
 "cat VERSION"
 "bump -patch -write"
 "cat VERSION"
 "bump -major -write"
 "cat VERSION"
 "bump -preview -write"
 "rm VERSION"
 "echo \"--- SCENARIO THREE ---\""
 "echo \"1.25\" > VERSION"
 "bump -fix"
 "bump -fix -write"
 "rm VERSION"
 "echo \"--- SCENARIO FOUR ---\""
 "echo \"v1.17.7-beta.6\" > VERSION"
 "bump -check -fix"
 "cat VERSION"
 "bump -check -fix -write"
 "cat VERSION"
 "rm VERSION"
 "echo \"--- SCENARIO FIVE ---\""
 "echo \"module testApp-${counterName}\" > go.mod"
 "echo \"\" >> go.mod"
 "echo \"go 1.24\" >> go.mod"
 "cat go.mod"
 "bump -in go.mod -fix"
 "bump -in go.mod -fix -write"
 "cat go.mod"
 "rm go.mod"
 "echo \"--- SCENARIO SIX ---\""
 "echo v1.0.0 > VERSION"
 "bump -json -check"
 "cat VERSION"
 "bump -json -beta"
 "cat VERSION"
 "bump -json -beta -write"
 "cat VERSION"
 "rm VERSION"
 "echo \"--- SCENARIO SEVEN: package.json ---\""
 "echo eyJuYW1lIjoidGVzdCIsInZlcnNpb24iOiIxLjIuMyJ9 | base64 -d | tee package.json > /dev/null"
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
  "echo \"--- SCENARIO EIGHT: go.mod ---\""
  "echo 'module myapp' > go.mod && echo 'go 1.21' >> go.mod"
  "bump -in go.mod -check"
  "bump -in go.mod -minor"
  "bump -in go.mod -minor -write"
  "grep 'go 1.22' go.mod"
  "rm go.mod"
  "echo \"--- SCENARIO NINE: Dockerfile ---\""
  "echo 'LABEL version=\"v3.2.1\"' > Dockerfile"
  "bump -in Dockerfile -check"
  "bump -in Dockerfile -patch"
  "bump -in Dockerfile -patch -write"
  "grep 'LABEL version=\"v3.2.2\"' Dockerfile"
  "rm Dockerfile"
  "echo \"--- SCENARIO TEN: Chart.yaml (Helm) ---\""
  "echo 'YXBpVmVyc2lvbjogdjIKbmFtZTogbXljaGFydAp2ZXJzaW9uOiAwLjEuMAo=' | base64 -d | tee Chart.yaml > /dev/null"
  "bump -in Chart.yaml -check"
  "bump -in Chart.yaml -patch"
  "bump -in Chart.yaml -patch -write"
  "grep 'version: 0.1.1' Chart.yaml"
  "rm Chart.yaml"
  "echo \"--- SCENARIO ELEVEN: pom.xml (Maven) ---\""
  "echo '<project><version>2.2.2</version></project>' > pom.xml"
  "bump -in pom.xml -check"
  "bump -in pom.xml -patch"
  "bump -in pom.xml -patch -write"
  "grep '<version>2.2.3</version>' pom.xml"
  "rm pom.xml"
  "echo \"--- SCENARIO TWELVE: Environment Variables ---\""
  "echo \"v5.5.5\" > VERSION"
  "BUMP_ALWAYS_WRITE=true bump -env"
  "BUMP_ALWAYS_WRITE=true bump -patch"
  "grep 'v5.5.6' VERSION"
  "BUMP_DEFAULT_INPUT=VERSION bump -minor"
  "BUMP_DEFAULT_INPUT=VERSION bump -minor -write"
  "grep 'v5.6.0' VERSION"
  "rm VERSION"
  "echo \"--- TEST THIRTEEN ---\""
  "bump -parse v1.2.3-alpha.4 -init"
  "cat VERSION"
  "bump -parse v2.3.4-alpha.5 -write"
  "cat VERSION"
  "rm VERSION"
  "bump -parse v3.4.5-alpha.6 -init"
)

for t in "${tests[@]}"; do
  run "${t}"
done

echo "All $(counter -name $counterName) tests PASS!"
