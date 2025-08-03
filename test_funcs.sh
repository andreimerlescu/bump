#!/bin/bash

[ "${IN_TEST_FILE}" != "1" ] && { echo "Cannot call this script outside of test.sh!"; exit 1; }

declare populated_package_json
populated_package_json="ewogICJuYW1lIjogIm15X3BhY2thZ2UiLAogICJkZXNjcmlwdGlvbiI6ICJtYWtlIHlvdXIgcGFja2FnZSBlYXNpZXIgdG8gZmluZCBvbiB0aGUgbnBtIHdlYnNpdGUiLAogICJ2ZXJzaW9uIjogIjEuMC4wIiwKICAic2NyaXB0cyI6IHsKICAgICJ0ZXN0IjogImVjaG8gXCJFcnJvcjogbm8gdGVzdCBzcGVjaWZpZWRcIiAmJiBleGl0IDEiCiAgfSwKICAicmVwb3NpdG9yeSI6IHsKICAgICJ0eXBlIjogImdpdCIsCiAgICAidXJsIjogImh0dHBzOi8vZ2l0aHViLmNvbS9tb25hdGhlb2N0b2NhdC9teV9wYWNrYWdlLmdpdCIKICB9LAogICJrZXl3b3JkcyI6IFtdLAogICJhdXRob3IiOiAiIiwKICAibGljZW5zZSI6ICJJU0MiLAogICJidWdzIjogewogICAgInVybCI6ICJodHRwczovL2dpdGh1Yi5jb20vbW9uYXRoZW9jdG9jYXQvbXlfcGFja2FnZS9pc3N1ZXMiCiAgfSwKICAiaG9tZXBhZ2UiOiAiaHR0cHM6Ly9naXRodWIuY29tL21vbmF0aGVvY3RvY2F0L215X3BhY2thZ2UiCn0="
export populated_package_json

declare empty_package_json
empty_package_json="eyJuYW1lIjoidGVzdCIsInZlcnNpb24iOiIxLjIuMyJ9"
export empty_package_json

declare empty_chart_yaml
empty_chart_yaml="YXBpVmVyc2lvbjogdjIKbmFtZTogbXljaGFydAp2ZXJzaW9uOiAwLjEuMAo="
export empty_chart_yaml

declare counterName

function safe_exit() {
    local msg="${1}"
    echo -e "${msg}"
    exit 1
}

function run() {
  local cmd="${1}"
  local -i testNo
  testNo=$(counter -name "${counterName}" -add || safe_exit "failed to increase counter")
  local -i exitCode
  local output

  # Capture both stdout and stderr together
  output=$(eval "${cmd}" 2>&1)
  exitCode=$?

  local prefix
  prefix="$(magenta "$(whoami)")@$(yellow bump.git):$(purple "$(basename "${0}")")"

  if (( exitCode == 0 )); then
    printf "%s ⚡ %s ⇒  %s\n" "${prefix}" "$(cyan "Test #${testNo}")" "$(green "${cmd}")"
    [[ -n "$output" ]] && printf "%s\n" "$output"
  else
    printf "%s ⚡ %s ⇒  %s\n" "${prefix}" "$(cyan "Test #${testNo}")" "$(red "${cmd}")"
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

function red() { if [ "${NO_COLOR}" == "false" ]; then echo -e "\033[0;31m${1}\033[0m"; else echo "${1}"; fi; }
function purple() { if [ "${NO_COLOR}" == "false" ]; then echo -e "\033[0;34m${1}\033[0m";  else echo "${1}"; fi;}
function green() { if [ "${NO_COLOR}" == "false" ]; then echo -e "\033[0;32m${1}\033[0m";  else echo "${1}"; fi;}
function magenta() { if [ "${NO_COLOR}" == "false" ]; then echo -e "\033[0;35m${1}\033[0m";  else echo "${1}"; fi;}
function yellow() { if [ "${NO_COLOR}" == "false" ]; then echo -e "\033[1;33m${1}\033[0m";  else echo "${1}"; fi;}
function cyan() { if [ "${NO_COLOR}" == "false" ]; then echo -e "\033[1;36m${1}\033[0m";  else echo "${1}"; fi;}

function cleanup() {
  unset scenario_01
  unset scenario_02
  unset scenario_03
  unset scenario_04
  unset scenario_05
  unset scenario_06
  unset scenario_07
  unset scenario_08
  unset scenario_09
  unset scenario_10
  unset scenario_11
  unset scenario_12
  unset scenario_13
  unset scenario_14
  unset tests
  unset populated_package_json
  unset empty_chart_yaml
  unset empty_package_json
  unset counterName
  unset COUNTER_DIR
  unset COUNTER_ALWAYS_YES
  unset COUNTER_USE_FORCE
  unset NO_COLOR
  unset IN_TEST_FILE
  cd ..
  rm -rf test-data
}
