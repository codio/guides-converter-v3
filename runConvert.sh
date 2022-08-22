#!/bin/bash

converterVersion=$1


curl "https://static-assets.codio.com/guides-converter-v3/guides-converter-v3-${converterVersion}" --output guides-converter-v
chmod +x ./guides-converter-v3
./guides-converter-v3
rm guides-converter-v3

exit $?
