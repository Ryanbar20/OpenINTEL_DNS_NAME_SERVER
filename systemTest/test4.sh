#!/bin/bash

# start the name-server
./../main/main 1 10 [::]:10000 >/dev/null &
pid=$!
sleep 3 

#clear output file
echo '' > out4.txt 


echo "============================================================" >> out4.txt
echo 'submit a query, wait and observe that the retry returns an answer' >> out4.txt
echo "============================================================" >> out4.txt
dig @127.0.0.1 -p 10000 "20231001.google.nu.history.openintel.nl" A >> out4.txt
while true; do
    sleep 2
    echo "performing a retry" >> out4.txt
    output=$(dig @127.0.0.1 -p 10000 "20231001.google.nu.history.openintel.nl" A)
    if ! grep -q 'HINFO "WAIT"' <<< "$output"; then
        echo "$output" >> out4.txt
        break
    fi
done
echo "============================================================" >> out4.txt
echo 'send a different query, wait and observe that the retry returns an answer' >> out4.txt
echo "============================================================" >> out4.txt
dig @127.0.0.1 -p 10000 "20231002.google.nu.history.openintel.nl" A >> out4.txt
while true; do
    sleep 2
    echo "performing a retry" >> out4.txt
    output=$(dig @127.0.0.1 -p 10000 "20231002.google.nu.history.openintel.nl" A)
    if ! grep -q 'HINFO "WAIT"' <<< "$output"; then
        echo "$output" >> out4.txt
        break
    fi
done
echo "============================================================" >> out4.txt
echo 'retry the first query, observe that it returns a wait-message again' >> out4.txt
echo "============================================================" >> out4.txt
dig @127.0.0.1 -p 10000 "20231001.google.nu.history.openintel.nl" A >> out4.txt


# kill the name server
kill $pid
