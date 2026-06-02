#!/bin/bash

# start the name-server
./../main/main 1 10 [::]:10000 &
pid=$!
sleep 3 

# submit a query, wait and observe that the retry returns an answer
dig @127.0.0.1 -p 10000 "20231001.google.nu.history.openintel.nl" A > out4.txt
sleep 20
dig @127.0.0.1 -p 10000 "20231001.google.nu.history.openintel.nl" A >> out4.txt

#send a different query, wait and observe that the retry returns an answer
dig @127.0.0.1 -p 10000 "20231002.google.nu.history.openintel.nl" A >> out4.txt
sleep 20
dig @127.0.0.1 -p 10000 "20231002.google.nu.history.openintel.nl" A >> out4.txt

#retry the first query, observe that it returns a 'wait' again
dig @127.0.0.1 -p 10000 "20231001.google.nu.history.openintel.nl" A >> out4.txt


# kill the name server
kill $pid 
