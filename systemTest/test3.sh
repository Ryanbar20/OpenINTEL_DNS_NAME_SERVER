#!/bin/bash

# start the name-server
./../main/main 10 10 [::]:10000 4GB >/dev/null &
pid=$!
sleep 3 

#clear output file
echo '' > out3.txt 


echo "============================================================" >> out3.txt
echo 'submit a query' >> out3.txt
echo "============================================================" >> out3.txt
dig @127.0.0.1 -p 10000 "20231001.google.nu.history.openintel.nl" A >> out3.txt

echo "============================================================" >> out3.txt
echo 'retry the query on 2 second intervals until it returns an answer' >> out3.txt
echo "============================================================" >> out3.txt
while true; do
    sleep 2
    output=$(dig @127.0.0.1 -p 10000 "20231001.google.nu.history.openintel.nl" A)
	echo "$output" >> out3.txt
    if ! grep -q 'HINFO "WAIT"' <<< "$output"; then
        break
    fi
done

echo "============================================================" >> out3.txt
echo 'retry the query on 2 second intervals and see that it keeps on giving the answer' >> out3.txt
echo "============================================================" >> out3.txt
for i in {1..5}
do 
	dig @127.0.0.1 -p 10000 "20231001.google.nu.history.openintel.nl" A >> out3.txt
	sleep 2 
done

# kill the name server
kill $pid
