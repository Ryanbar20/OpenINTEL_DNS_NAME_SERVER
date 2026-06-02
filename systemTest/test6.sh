#!/bin/bash

# start the name-server
./../main/main 10 10 [::]:10000 &
pid=$!
sleep 3 

# add 2 distinct IPs to dig from
sudo ip addr add 127.0.0.2/8 dev lo
sudo ip addr add 127.0.0.3/8 dev lo

# submit a query
dig @127.0.0.1 -p 10000 -b 127.0.0.1 "20231001.google.nu.history.openintel.nl" A > out6.txt
# send the query again but with a different IP. Observe that a 'wait' message with the same queue position as the previous query is returned
dig @127.0.0.1 -p 10000 -b 127.0.0.3 "20231001.google.nu.history.openintel.nl" A >> out6.txt
# send a different query with a different IP, Observe that a 'wait' message is returned with another queue position
dig @127.0.0.1 -p 10000 -b 127.0.0.3 "20231002.google.nu.history.openintel.nl" A >> out6.txt

# kill the name server
kill $pid 

sudo ip addr del 127.0.0.2/8 dev lo
sudo ip addr del 127.0.0.3/8 dev lo