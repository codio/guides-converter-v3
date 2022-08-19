#!/bin/bash

CONVERTER_VERSION="72ecc41ef93e390bf6d9287eb13b9809e4f82c60"

curl "https://static-assets.codio.com/guides-converter-v3/guides-converter-v3-${CONVERTER_VERSION}" --output guides-converter-v3
chmod +x ./guides-converter-v3
./guides-converter-v3
rm guides-converter-v3
