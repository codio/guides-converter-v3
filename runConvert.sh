#!/bin/bash

converterVersion=$1

run_with_failover () {
  "$@"
  ret=$?
  if [[ $ret -ne 0 ]]
  then
      rm guides-converter-v3
      exit $ret
  fi
}

curl "https://static-assets.codio.com/guides-converter-v3/guides-converter-v3-${converterVersion}" --output guides-converter-v3

run_with_failover chmod +x ./guides-converter-v3

run_with_failover ./guides-converter-v3

rm guides-converter-v3
