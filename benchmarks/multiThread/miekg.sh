#!/bin/bash


# clear the output files

out="MiekgOut.txt"
debug="MiekgDebug.txt"

echo '' > "$out"
echo '' > "$debug"




for i in {1..30} ;
do
    # start a new name-server
    ./miekg >/dev/null &
    pid=$!
    sleep 3 
    
    echo "============= DNSPERF ===========" >> "$out"
    dnsperf -m udp -s 127.0.0.1 -p 10000 -d ./dnsperf-queries.txt -c 4 -q 200 -l 30 >> "$out"


    # kill the name server
    kill $pid

done

