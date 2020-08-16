#! /bin/bash

openssl genrsa -out ../certs/server.key 2048
openssl rsa -in ../certs/server.key -pubout > ../certs/publickey.pub
openssl req -new -x509 -sha256 -key ../certs/server.key -out ../certs/server.crt -days 5475 -subj "/C=RO/ST=Bucharest/L=Bucharest/O=CryptoCube SRL/OU=CORE/CN=dsus.it"

