#!bin/bash

domain="qsocks.org"
email="admin@qsocks.org"

echo "make server cert"
openssl req -new -nodes -x509 -out ./certs/server.pem -keyout ./certs/server.key -days 3650 -subj "/C=DE/ST=NRW/L=Earth/O=Random Company/OU=IT/CN=$domain/emailAddress=$email"
echo "make client cert"
openssl req -new -nodes -x509 -out ./certs/client.pem -keyout ./certs/client.key -days 3650 -subj "/C=DE/ST=NRW/L=Earth/O=Random Company/OU=IT/CN=$domain/emailAddress=$email"

