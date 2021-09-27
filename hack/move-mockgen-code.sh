#!/usr/bin/env bash

TARGETDIR="internal/mockgen/native/runners/runner"
TARGETDIR2="internal/mockgen/native/runnertemplate"

# get golang/mock submodule
git submodule update --init

cd submodules/gomock

for VAL in `seq 32 -1 1`
do
  mkdir "../../${TARGETDIR}${VAL}"
  cp -r mockgen/* "../../${TARGETDIR}${VAL}"
  grep -lr "package" "../../${TARGETDIR}${VAL}/"* | xargs sed -i "" "1,20s/package main/package runner${VAL}/g"
  grep -lr "flag.String" "../../${TARGETDIR}${VAL}/"* | xargs sed -i "" "s/= flag\.String.*/\*string/g"
  grep -lr "flag.Bool" "../../${TARGETDIR}${VAL}/"* | xargs sed -i "" "s/= flag\.Bool.*/\*bool/g"

  grep -lr "var output" "../../${TARGETDIR}${VAL}/"* | xargs sed -i "" "s/var output \*string/var output = flag.String(\"output\", \"\", \"The output file name, or empty to use stdout\.\")/g"

  find ../../${TARGETDIR}${VAL}/ -print | grep --regex '.*\.go' | xargs goimports -w
done

cp -r mockgen/* "../../${TARGETDIR2}"
grep -lr "package" "../../${TARGETDIR2}/"* | xargs sed -i "" "1,20s/package main/package runnertemplate/g"
grep -lr "flag.String" "../../${TARGETDIR2}/"* | xargs sed -i "" "s/= flag\.String.*/\*string/g"
grep -lr "flag.Bool" "../../${TARGETDIR2}/"* | xargs sed -i "" "s/= flag\.Bool.*/\*bool/g"

grep -lr "var output" "../../${TARGETDIR2}/"* | xargs sed -i "" "s/var output \*string/var output = flag.String(\"output\", \"\", \"The output file name, or empty to use stdout\.\")/g"

find ../../${TARGETDIR2}/ -print | grep --regex '.*\.go' | xargs goimports -w