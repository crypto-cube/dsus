#! /bin/bash

if id "dsus" &>/dev/null; then
    echo 'user exists skipping ...'
else
    useradd -m dsus
fi

chmod 755 /etc/dsus
chmod 644 /etc/dsus/* 