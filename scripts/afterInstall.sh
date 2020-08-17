#! /bin/bash

if id "dsus" &>/dev/null; then
    echo 'user exists skipping ...'
else
    useradd -m dsus
fi

chmod 755 -R /etc/dsus/ && echo "Permissions set!"