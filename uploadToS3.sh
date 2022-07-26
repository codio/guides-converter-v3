#!/bin/bash
set -xe

s3Key=$1
s3Secret=$2

uploadFile () {
    fName=$1
    
    file=./$fName
    bucket=codio-assets
    resource="/${bucket}/guides-converter-v3/${fName}"
    contentType="application/octet-stream"
    dateValue=`date -R`
    stringToSign="PUT\n\n${contentType}\n${dateValue}\n${resource}"
    signature=`echo -en ${stringToSign} | openssl sha1 -hmac ${s3Secret} -binary | base64`
    curl -X PUT -T "${file}" \
      -H "Host: ${bucket}.s3.amazonaws.com" \
      -H "Date: ${dateValue}" \
      -H "Content-Type: ${contentType}" \
      -H "Authorization: AWS ${s3Key}:${signature}" \
      https://${bucket}.s3.amazonaws.com/guides-converter-v3/${fName}
}

cd dist
for file in *; do uploadFile $file; done
