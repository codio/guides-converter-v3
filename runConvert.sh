#!/bin/bash

converterVersion=$1
exit_code=0


check_exit_code () {
  code=$1
  if [ $code -eq 1 ]
  then
    exit_code=1
  fi
}

curl "https://static-assets.codio.com/guides-converter-v3/guides-converter-v3-${converterVersion}" --output guides-converter-v3
check_exit_code $?

chmod +x ./guides-converter-v3
check_exit_code $?

./guides-converter-v3
check_exit_code $?

rm guides-converter-v3
check_exit_code $?
exit $exit_code
