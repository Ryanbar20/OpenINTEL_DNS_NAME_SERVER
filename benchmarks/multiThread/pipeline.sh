#!/bin/bash

if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root"
   echo "This is because it needs to create an additional IP loopback address to locally send dns queries from"
   echo "This address will be deleted afterwards"
   exit 1
fi

# add 2 distinct IPs to dig from
sudo ip addr add 127.0.0.2/8 dev lo
sudo ip addr add 127.0.0.3/8 dev lo

# clear the output files
echo '' > out.txt
echo '' > debug.txt

# start the name-server
./../../main/main 10 10 [::]:10000 >/dev/null &
pid=$!
sleep 3 


for i in {1..30} ;
do
    {
        # seed the cache
    dig @127.0.0.1 -p 10000 "20201001.google.nu.history.openintel.nl" A
    while true; do
        sleep 2
        echo "performing a retry"
        output=$(dig @127.0.0.1 -p 10000 "20201001.google.nu.history.openintel.nl" A)
        if ! grep -q 'HINFO "WAIT"' <<< "$output"; then
            break
        fi
    done

    dig @127.0.0.1 -p 10000 "20201001.google.nu.history.openintel.nl" AAAA
    while true; do
        sleep 2
        echo "performing a retry"
        output=$(dig @127.0.0.1 -p 10000 "20201001.google.nu.history.openintel.nl" AAAA)
        if ! grep -q 'HINFO "WAIT"' <<< "$output"; then
            break
        fi
    done

    dig @127.0.0.1 -p 10000 "20201001.google.nu.history.openintel.nl" NS
    while true; do
        sleep 2
        echo "performing a retry"
        output=$(dig @127.0.0.1 -p 10000 "20201001.google.nu.history.openintel.nl" NS)
        if ! grep -q 'HINFO "WAIT"' <<< "$output"; then
            break
        fi
    done
    } > /dev/null 2>&1

    #confirm the cache
    checkOutput=$(
        echo "============= CONFIRM THE CACHE ENTRIES ==========="
        {
            output=$(dig @127.0.0.1 -p 10000 "20201001.google.nu.history.openintel.nl" A)
            if grep -q 'HINFO "WAIT"' <<< "$output"; then
                echo "ERR: Cache was not correctly filled"
            fi
            echo "$output"
        }
        {
            output=$(dig @127.0.0.1 -p 10000 "20201001.google.nu.history.openintel.nl" AAAA)
            if grep -q 'HINFO "WAIT"' <<< "$output"; then
                echo "ERR: Cache was not correctly filled"
            fi
            echo "$output"
        }
        {
            output=$(dig @127.0.0.1 -p 10000 "20201001.google.nu.history.openintel.nl" NS)
            if grep -q 'HINFO "WAIT"' <<< "$output"; then
                echo "ERR: Cache was not correctly filled"
            fi
            echo "$output"
        }



        #fill and confirm the queue
        echo "============= FILL THE QUEUE ==========="
        {
            output=$(dig @127.0.0.1 -p 10000 -b 127.0.0.1 "20231001.google.fr.history.openintel.nl" A)
            if ! grep -q '"You are in queue position 0"' <<< "$output"; then
                echo "ERR: Queue was not correctly filled"
            fi
            echo "$output"
        }
        {
            output=$(dig @127.0.0.1 -p 10000 -b 127.0.0.2 "20231001.google.fr.history.openintel.nl" AAAA)
            if ! grep -q '"You are in queue position 1"' <<< "$output"; then
                echo "ERR: Queue was not correctly filled"
            fi
            echo "$output"
        }
        {
            output=$(dig @127.0.0.1 -p 10000 -b 127.0.0.3 "20231002.google.fr.history.openintel.nl" A)
            if ! grep -q '"You are in queue position 2"' <<< "$output"; then
                echo "ERR: Queue was not correctly filled"
            fi
            echo "$output"
        }
    )

    echo "$checkOutput" >> debug.txt
    if grep -q '"ERR: "' <<< "$checkOutput"; then
        echo 'INVALID TEST' >> debug.txt
    else 
        #now dnsperf
        echo "============= DNSPERF ===========" >> out.txt
        dnsperf -m udp -s 127.0.0.1 -p 10000 -d ./dnsperf-queries.txt -c 4 -q 200 -l 30 >> out.txt


        echo "============= CONFIRM QUEUE ENTRIES (at least 1 should still be in the queue so just check the last query) ===========" >> out.txt
        dig @127.0.0.1 -p 10000 "20231002.google.fr.history.openintel.nl" A >> out.txt
    fi

    # kill the name server
    kill $pid

done




sudo ip addr del 127.0.0.2/8 dev lo
sudo ip addr del 127.0.0.3/8 dev lo

