#!/bin/bash

# start the name-server
./../main/main 10 2 [::]:10000 &
pid=$!
sleep 3 

# add 2 distinct IPs to dig from
sudo ip addr add 127.0.0.2/8 dev lo
sudo ip addr add 127.0.0.3/8 dev lo

# submit three queries. Observe that the last gets answered with a limit message as the queue is full
dig @127.0.0.1 -p 10000 -b 127.0.0.1 "20231001.google.nu.history.openintel.nl" A > out7.txt
dig @127.0.0.1 -p 10000 -b 127.0.0.2 "20231002.google.nu.history.openintel.nl" A >> out7.txt
dig @127.0.0.1 -p 10000 -b 127.0.0.3 "20231003.google.nu.history.openintel.nl" A >> out7.txt

# kill the name server
kill $pid 

sudo ip addr del 127.0.0.2/8 dev lo
sudo ip addr del 127.0.0.3/8 dev lo