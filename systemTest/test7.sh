#!/bin/bash

# start the name-server
./../main/main 10 2 [::]:10000 >/dev/null &
pid=$!
sleep 3 


#clear output file
echo '' > out7.txt 


echo "============================================================" >> out7.txt
echo 'submit a query and wait. Observe that the retry returns an answer' >> out7.txt
echo "============================================================" >> out7.txt
dig @127.0.0.1 -p 10000  "20231001.google.nu.history.openintel.nl" A >> out7.txt
sleep 20
dig @127.0.0.1 -p 10000  "20231001.google.nu.history.openintel.nl" A >> out7.txt

echo "============================================================" >> out7.txt
echo 'submit a different query' >> out7.txt
echo "============================================================" >> out7.txt
dig @127.0.0.1 -p 10000  "20231002.google.nu.history.openintel.nl" A >> out7.txt

echo "============================================================" >> out7.txt
echo 'resubmit the first query. Observe that it returns an answer' >> out7.txt
echo "============================================================" >> out7.txt
dig @127.0.0.1 -p 10000  "20231001.google.nu.history.openintel.nl" A >> out7.txt

echo "============================================================" >> out7.txt
echo 'resubmit the second query. Observe that it is still in the queue' >> out7.txt
echo "============================================================" >> out7.txt
dig @127.0.0.1 -p 10000  "20231002.google.nu.history.openintel.nl" A >> out7.txt
sleep 20

echo "============================================================" >> out7.txt
echo 'resubmit the second query after waiting 20 seconds. Observe that it is now done processing' >> out7.txt
echo "============================================================" >> out7.txt
dig @127.0.0.1 -p 10000  "20231002.google.nu.history.openintel.nl" A >> out7.txt

# kill the name server
kill $pid
