#! /bin/bash

set -x

FILE=$1
SERVER=https://127.0.0.1:8787/upload

[ -z "$FILE" ] && echo "No file specified" && exit 1;

COMPRESSED=$(mktemp /tmp/dsus-upload-XXXXXX.zst)
zstd -T0 -c "$FILE" > "$COMPRESSED"

openssl dgst -sha256 -sign ../certs/server.key -out /tmp/signature1 "$COMPRESSED"
curl -F "executable=@$COMPRESSED" -F "signature=@/tmp/signature1" "$SERVER" --insecure

rm -f "$COMPRESSED"