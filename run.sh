#!/usr/bin/env bash

grep "127.0.0.1 localhost localhost1 localhost2" /etc/hosts &> /dev/null

if [[ $? -ne 0 ]]; then
    echo -e "\n127.0.0.1 localhost localhost1 localhost2" >> /etc/hosts
fi

echo "$(tput setaf 2)----------------------------------------------------------------------$(tput sgr 0)"

echo "$(tput setaf 2)[INFO]$(tput sgr 0) Removing previous generated keys"
rm proxy/*.key && rm proxy/*.crt && rm proxy/*.csr && rm proxy/*.log 2> /dev/null


echo "$(tput setaf 2)[INFO]$(tput sgr 0) Creating self-signed certificate for localhost1 on dir proxy..."
openssl req  -new  -newkey rsa:2048  -nodes  -keyout proxy/localhost1.key -subj "/CN=localhost1/emailAddress=admin@mail/C=BR/ST=RJ/L=Rio de Janeiro/O=Empty/OU=Empty"  -out proxy/localhost1.csr &> proxy/certificate1.log
openssl  x509  -req  -days 365  -in proxy/localhost1.csr  -signkey proxy/localhost1.key  -out proxy/localhost1.crt &>> proxy/certificate1.log

echo "$(tput setaf 2)[INFO]$(tput sgr 0) Creating self-signed certificate for localhost2 on dir proxy..."
openssl req  -new  -newkey rsa:2048  -nodes  -keyout proxy/localhost2.key -subj "/CN=localhost2/emailAddress=admin@mail/C=BR/ST=RJ/L=Rio de Janeiro/O=Empty/OU=Empty"  -out proxy/localhost2.csr &> proxy/certificate2.log
openssl  x509  -req  -days 365  -in proxy/localhost2.csr  -signkey proxy/localhost2.key  -out proxy/localhost2.crt &>> proxy/certificate2.log

echo "$(tput setaf 2)[INFO]$(tput sgr 0) Creating log file for operation above on dir proxy..."
echo "$(tput setaf 2)----------------------------------------------------------------------$(tput sgr 0)"


docker-compose up