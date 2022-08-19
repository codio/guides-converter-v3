#!/bin/bash

CONVERTER_VERSION="8e36af2f397cda139864f9117d2f3ef69ddc0bea"

curl "https://static-assets.codio.com/guides-converter-v3/guides-converter-v3-${CONVERTER_VERSION}" --output guides-converter-v3
chmod +x ./guides-converter-v3
./guides-converter-v3
rm guides-converter-v3
