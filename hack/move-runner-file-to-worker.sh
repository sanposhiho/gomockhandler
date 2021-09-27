#!/usr/bin/env bash

TARGETDIR="internal/mockgen/native"

for VAL in `seq 32 -1 1`
do
  cp "${TARGETDIR}/runnertemplate/runner.go" "${TARGETDIR}/runners/runner${VAL}/"
  grep -lr "package" "${TARGETDIR}/runners/runner${VAL}/runner.go" | xargs sed -i "" "s/package runnertemplate/package runner${VAL}/g"
done
