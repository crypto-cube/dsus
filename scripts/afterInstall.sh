#! /bin/bash

if id "dsus" &>/dev/null; then
    echo 'user exists skipping ...'
else
    useradd -m dsus
fi

if [ ! -f "/etc/dsus/certs/publickey.pub" ]; then
    echo "Generating keys ..."
    openssl genrsa -out /etc/dsus/certs/server.key 2048
    openssl rsa -in /etc/dsus/certs/server.key -pubout > /etc/dsus/certs/publickey.pub
    openssl req -new -x509 -sha256 -key /etc/dsus/certs/server.key -out /etc/dsus/certs/server.crt -days 5475 -subj "/C=RO/ST=Bucharest/L=Bucharest/O=CryptoCube SRL/OU=CORE/CN=dsus.it"
fi

chmod 755 -R /etc/dsus/ && echo "Permissions set!"