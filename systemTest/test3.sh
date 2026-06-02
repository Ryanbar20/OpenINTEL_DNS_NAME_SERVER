#!/bin/bash

# start the name-server
./../main/main 10 10 [::]:10000 >/dev/null &
pid=$!
sleep 3 

#clear output file
echo '' > out3.txt 


echo "============================================================" >> out3.txt
echo 'submit a query' >> out3.txt
echo "============================================================" >> out3.txt
dig @127.0.0.1 -p 10000 "20231001.google.nu.history.openintel.nl" A >> out3.txt

echo "============================================================" >> out3.txt
echo 'retry the query 10 times on 2 second intervals' >> out3.txt
echo 'The first queries will be answered by a wait-message, the last by an answer' >> out3.txt
echo "============================================================" >> out3.txt
for i in {1..10}
do 
	dig @127.0.0.1 -p 10000 "20231001.google.nu.history.openintel.nl" A >> out3.txt
	sleep 2 
done

# kill the name server
kill $pid
