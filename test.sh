#!/bin/bash

export IN_TEST_FILE=1

[ ! -f test_funcs.sh ] && { echo "Missing test_funcs.sh script!"; exit 1; }
source test_funcs.sh

[ ! -f test_scenarios.sh ] && safe_exit "Missing test_scenario.sh script!"
source test_scenarios.sh

declare NO_COLOR
NO_COLOR="${NO_COLOR:-false}"
export NO_COLOR

go install github.com/andreimerlescu/counter@latest || safe_exit "failed to install counter"

[ ! -d .counters ] && { mkdir -p .counters || safe_exit "failed to mkdir .counters"; }

COUNTER_DIR=$(realpath .counters)

export COUNTER_USE_FORCE=1
export COUNTER_ALWAYS_YES=1

declare COUNTER_DIR
export COUNTER_DIR

declare counterName
counterName="bump-test-passes"
export counterName
counter -name $counterName -reset -yes 1> /dev/null || safe_exit "failed to reset counter"

echo "Preparing test env..."
mkdir -p test-data || safe_exit "failed to mkdir test-data"
go build -o ./test-data/bump . || safe_exit "failed to build bump"
cd test-data || safe_exit "failed to cd test-data"
chmod +x bump || safe_exit "failed to chmod bump"

declare -a tests=(
 "${scenario_01[@]}"
 "${scenario_02[@]}"
 "${scenario_03[@]}"
 "${scenario_04[@]}"
 "${scenario_05[@]}"
 "${scenario_06[@]}"
 "${scenario_07[@]}"
 "${scenario_08[@]}"
 "${scenario_09[@]}"
 "${scenario_10[@]}"
 "${scenario_11[@]}"
 "${scenario_12[@]}"
 "${scenario_13[@]}"
 "${scenario_14[@]}"
)

for t in "${tests[@]}"; do
  run "${t}"
done

echo "All $(counter -name $counterName) tests PASS!"

cleanup
