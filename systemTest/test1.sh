#!/bin/bash

# start the name-server
./../main/main 10 10 [::]:10000 4GB >/dev/null &
pid=$!
sleep 3 

#clear output file
echo '' > out1.txt 

echo "============================================================" >> out1.txt
echo 'submit an A query and see the retry give the answer' >> out1.txt
echo "============================================================" >> out1.txt
dig @127.0.0.1 -p 10000 "20201001.google.nu.history.openintel.nl" A >> out1.txt
while true; do
    sleep 2
    echo "performing a retry" >> out1.txt
    output=$(dig @127.0.0.1 -p 10000 "20201001.google.nu.history.openintel.nl" A)
    if ! grep -q 'HINFO "WAIT"' <<< "$output"; then
        echo "$output" >> out1.txt
        break
    fi
done

echo "============================================================" >> out1.txt
echo 'submit an AAAA query and see the retry give the answer' >> out1.txt
echo "============================================================" >> out1.txt
dig @127.0.0.1 -p 10000 "20201001.google.nu.history.openintel.nl" AAAA >> out1.txt
while true; do
    sleep 2
    echo "performing a retry" >> out1.txt
    output=$(dig @127.0.0.1 -p 10000 "20201001.google.nu.history.openintel.nl" AAAA)
    if ! grep -q 'HINFO "WAIT"' <<< "$output"; then
        echo "$output" >> out1.txt
        break
    fi
done

echo "============================================================" >> out1.txt
echo 'submit an NS query and see the retry give the answer' >> out1.txt
echo "============================================================" >> out1.txt
dig @127.0.0.1 -p 10000 "20201001.google.nu.history.openintel.nl" NS >> out1.txt
while true; do
    sleep 2
    echo "performing a retry" >> out1.txt
    output=$(dig @127.0.0.1 -p 10000 "20201001.google.nu.history.openintel.nl" NS)
    if ! grep -q 'HINFO "WAIT"' <<< "$output"; then
        echo "$output" >> out1.txt
        break
    fi
done

echo "============================================================" >> out1.txt
echo 'submit an MX query and see the retry give the answer' >> out1.txt
echo "============================================================" >> out1.txt
dig @127.0.0.1 -p 10000 "20201001.google.nu.history.openintel.nl" MX >> out1.txt
while true; do
    sleep 2
    echo "performing a retry" >> out1.txt
    output=$(dig @127.0.0.1 -p 10000 "20201001.google.nu.history.openintel.nl" MX)
    if ! grep -q 'HINFO "WAIT"' <<< "$output"; then
        echo "$output" >> out1.txt
        break
    fi
done

echo "============================================================" >> out1.txt
echo 'submit a TXT query and see the retry give the answer' >> out1.txt
echo "============================================================" >> out1.txt
dig @127.0.0.1 -p 10000 "20201001.google.nu.history.openintel.nl" TXT >> out1.txt
while true; do
    sleep 2
    echo "performing a retry" >> out1.txt
    output=$(dig @127.0.0.1 -p 10000 "20201001.google.nu.history.openintel.nl" TXT)
    if ! grep -q 'HINFO "WAIT"' <<< "$output"; then
        echo "$output" >> out1.txt
        break
    fi
done
# kill the name server
kill $pid
