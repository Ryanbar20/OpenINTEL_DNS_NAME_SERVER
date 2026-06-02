#!/bin/bash

# start the name-server
./../main/main 10 2 [::]:10000 &
pid=$!
sleep 3 


# submit a query and wait. Observe that the retry returns an answer
dig @127.0.0.1 -p 10000  "20231001.google.nu.history.openintel.nl" A > out8.txt
sleep 20
dig @127.0.0.1 -p 10000  "20231001.google.nu.history.openintel.nl" A >> out8.txt

# submit a different query
dig @127.0.0.1 -p 10000  "20231002.google.nu.history.openintel.nl" A >> out8.txt
# resubmit the first query. Observe that it returns an answer
dig @127.0.0.1 -p 10000  "20231001.google.nu.history.openintel.nl" A >> out8.txt
# resubmit the second query. Observe that it is still in the queue
dig @127.0.0.1 -p 10000  "20231002.google.nu.history.openintel.nl" A >> out8.txt
sleep 20
# resubmit the second query. Observe that it is now done processing
dig @127.0.0.1 -p 10000  "20231002.google.nu.history.openintel.nl" A >> out8.txt

# kill the name server
kill $pid 
