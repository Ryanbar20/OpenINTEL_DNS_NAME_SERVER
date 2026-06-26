#!/bin/bash

if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root"
   echo "This is because it needs to create an additional IP loopback address to locally send dns queries from"
   echo "This address will be deleted afterwards"
   exit 1
fi

# start the name-server
./../main/main 10 2 [::]:10000 4GB >/dev/null &
pid=$!
sleep 3 

#clear output file
echo '' > out5.txt 

# add 2 distinct IPs to dig from
sudo ip addr add 127.0.0.2/8 dev lo
sudo ip addr add 127.0.0.3/8 dev lo

echo "============================================================" >> out5.txt
echo 'submit three queries with different IPs. Observe that the last gets answered with a limit message as the queue is full' >> out5.txt
echo "============================================================" >> out5.txt
dig @127.0.0.1 -p 10000 -b 127.0.0.1 "20231001.google.nu.history.openintel.nl" A >> out5.txt
dig @127.0.0.1 -p 10000 -b 127.0.0.2 "20231002.google.nu.history.openintel.nl" A >> out5.txt
dig @127.0.0.1 -p 10000 -b 127.0.0.3 "20231003.google.nu.history.openintel.nl" A >> out5.txt


echo "============================================================" >> out5.txt
echo 'Send a different query (different from all previous queries) from the first IP. Observe that it gets answered with a personal limit' >> out5.txt
echo "============================================================" >> out5.txt
dig @127.0.0.1 -p 10000 -b 127.0.0.1 "20201002.google.nu.history.openintel.nl" A >> out5.txt


# kill the name server
kill $pid 

sudo ip addr del 127.0.0.2/8 dev lo
sudo ip addr del 127.0.0.3/8 dev lo