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
while true; do
    sleep 2
    echo "performing a retry" >> out7.txt
    output=$(dig @127.0.0.1 -p 10000  "20231001.google.nu.history.openintel.nl" A)
    if ! grep -q 'HINFO "WAIT"' <<< "$output"; then
        echo "$output" >> out7.txt
        break
    fi
done
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

echo "============================================================" >> out7.txt
echo 'retry the second query. Observe that it is done processing at the end' >> out7.txt
echo "============================================================" >> out7.txt
while true; do
    sleep 2
    echo "performing a retry" >> out7.txt
    output=$(dig @127.0.0.1 -p 10000  "20231002.google.nu.history.openintel.nl" A)
    if ! grep -q 'HINFO "WAIT"' <<< "$output"; then
        echo "$output" >> out7.txt
        break
    fi
done

# kill the name server
kill $pid
