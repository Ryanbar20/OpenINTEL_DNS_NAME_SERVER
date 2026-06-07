#!/bin/bash

# clear the output files
out="RegularPipelineOut.txt"
debug="RegularPipelineDebug.txt"

echo '' > "$out"
echo '' > "$debug"




for i in {1..30} ;
do
    # start a new name-server
    ./../../main/main 10 10 [::]:10000 >/dev/null &
    pid=$!
    sleep 3 
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
    )

    echo "$checkOutput" >> "$debug"
    if grep -q '"ERR: "' <<< "$checkOutput"; then
        echo 'INVALID TEST' >> "$debug"
    else 
        #now dnsperf
        echo "============= DNSPERF ===========" >> "$out"
        dnsperf -m udp -s 127.0.0.1 -p 10000 -d ./dnsperf-queries.txt -c 4 -q 200 -l 30 >> "$out"
    fi

    # kill the name server
    kill $pid

done
