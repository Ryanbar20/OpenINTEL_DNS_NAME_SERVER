#!/bin/bash

# start the name-server
./../main/main 10 10 [::]:10000 &
pid=$!
sleep 3 
# submit a query
dig @127.0.0.1 -p 10000 "20231001.google.nu.history.openintel.nl" A > out3.txt
# retry the query 10 times on 2 second intervals
# The last query/queries will PROBABLY be answered with an answer
for i in {1..10}
do 
	dig @127.0.0.1 -p 10000 "20231001.google.nu.history.openintel.nl" A >> out3.txt
	sleep 2 
done

# kill the name server
kill $pid 
