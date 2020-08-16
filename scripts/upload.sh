#! /bin/bash

set -x

FILE=$1
SERVER=https://127.0.0.1:8787/upload

[ -z "$FILE" ] && echo "No file specified" && exit 1;

openssl dgst -sha256 -sign ../certs/server.key -out /tmp/signature1 "$FILE"
curl -F "executable=@$FILE" -F "signature=@/tmp/signature1" "$SERVER" --insecure