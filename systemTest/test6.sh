#!/bin/bash

if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root"
   echo "This is because it needs to create an additional IP loopback address to locally send dns queries from"
   echo "This address will be deleted afterwards"
   exit 1
fi


# start the name-server
./../main/main 10 10 [::]:10000 4GB >/dev/null &
pid=$!
sleep 3 

#clear output file
echo '' > out6.txt 

# add an IP to dig from
sudo ip addr add 127.0.0.2/8 dev lo


echo "============================================================" >> out6.txt
echo 'submit a query' >> out6.txt
echo "============================================================" >> out6.txt
dig @127.0.0.1 -p 10000 -b 127.0.0.1 "20231001.google.nu.history.openintel.nl" A >> out6.txt

echo "============================================================" >> out6.txt
echo 'send the query again but with a different IP. Observe that a 'wait' message with the same queue position as the previous query is returned' >> out6.txt
echo "============================================================" >> out6.txt
dig @127.0.0.1 -p 10000 -b 127.0.0.2 "20231001.google.nu.history.openintel.nl" A >> out6.txt

echo "============================================================" >> out6.txt
echo 'send a different query with the same IP as the previous. Observe that a 'wait' message is returned with another queue position' >> out6.txt
echo "============================================================" >> out6.txt
dig @127.0.0.1 -p 10000 -b 127.0.0.2 "20231002.google.nu.history.openintel.nl" A >> out6.txt

# kill the name server
kill $pid

sudo ip addr del 127.0.0.2/8 dev lo