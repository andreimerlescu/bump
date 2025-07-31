#!/bin/bash

function safe_exit() {
    local msg="${1}"
    echo "${msg}"
    exit 1
}

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

function run() {
  local cmd="${1}"
  local -i testNo
  testNo=$(counter -name $counterName -add || safe_exit "failed to increase counter")
  local out
  local -i exitCode
  touch out.txt
  { eval "${cmd}" > out.txt; } 2>&1
  exitCode=$?
  local check=""
  if (( exitCode == 0 )); then
    check="SUCCESS"
  else
    check="FAILED"
    rm out.txt
    safe_exit "failed to eval ${cmd}"
  fi
  printf "run.%02d(%s)\n    stdout: %s\n" "${testNo}" "${cmd}" "$(cat out.txt)"
  rm out.txt

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
 "bump -gomod -in go.mod -fix"
 "cat go.mod"
 "bump -gomod -in go.mod -fix -write"
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
)

for t in "${tests[@]}"; do
  run "${t}"
done

echo "All $(counter -name $counterName) tests PASS!"