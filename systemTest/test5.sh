#!/bin/bash

# start the name-server
./../main/main 10 10 [::]:10000 &
pid=$!
sleep 3 

# submit two queries, observe that the second gets answered with a limit messsage
dig @127.0.0.1 -p 10000 "20231001.google.nu.history.openintel.nl" A > out5.txt
dig @127.0.0.1 -p 10000 "20231002.google.nu.history.openintel.nl" A >> out5.txt


# kill the name server
kill $pid 
