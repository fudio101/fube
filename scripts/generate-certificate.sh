#!/usr/bin/bash

cd ../cert/ || exit

ip=$(hostname -I | awk '{print $1}')
echo "$ip"

openssl genrsa -out myCA.key 2048

openssl req -x509 -new -key myCA.key -out myCA.cer -days 730 -subj /CN=$ip

openssl genrsa -out mycert1.key 2048

openssl req -new -out mycert1.req -key mycert1.key -subj /CN=$ip

openssl x509 -req -in mycert1.req -out mycert1.cer -CAkey myCA.key -CA myCA.cer -days 365 -CAcreateserial -CAserial serial

cd ../