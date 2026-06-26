#!/bin/bash




start="2023-01-01"


# start the name-server
/usr/bin/time -v ./../main/main 10 10 [::]:10000 1GB >> out8.txt 2>&1 &
sleep 3 

#clear output file
echo '' > out8.txt 

echo "============================================================" >> out8.txt
echo 'submit n queries.' >> out8.txt
echo "============================================================" >> out8.txt

for ((i=1; i<=$1; i++));
do
    query=$(date -d "$start + $i days" +"%Y%m%d.google.fr.history.openintel.nl.")
    while true; do
        sleep 2
        output=$(dig @127.0.0.1 -p 10000 "$query" A)
        if ! grep -q 'HINFO "WAIT"' <<< "$output"; then
            echo "$output" >> out8.txt
            break
        fi
    done;
done;
echo "============================================================" >> out8.txt
echo 'Observe the name-server process memory usage (should be around 1 GB).' >> out8.txt
echo "============================================================" >> out8.txt
pkill main  


echo "============================================================" >> out8.txt
echo 'Start another name server' >> out8.txt
echo "============================================================" >> out8.txt
# start the name-server
/usr/bin/time -v ./../main/main 10 10 [::]:10000 500MB >> out8.txt 2>&1 &
sleep 3 

#clear output file
echo '' > out8.txt 

echo "============================================================" >> out8.txt
echo 'submit n queries.' >> out8.txt
echo "============================================================" >> out8.txt

for ((i=1; i<=$1; i++));
do
    query=$(date -d "$start + $i days" +"%Y%m%d.google.fr.history.openintel.nl.")
    while true; do
        sleep 2
        output=$(dig @127.0.0.1 -p 10000 "$query" A)
        if ! grep -q 'HINFO "WAIT"' <<< "$output"; then
            echo "$output" >> out8.txt
            break
        fi
    done;
done;
echo "============================================================" >> out8.txt
echo 'Observe the name-server process memory usage (should be around 500MB).' >> out8.txt
echo "============================================================" >> out8.txt
pkill main  