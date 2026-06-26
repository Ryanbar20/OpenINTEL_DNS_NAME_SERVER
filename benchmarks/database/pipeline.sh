#!/bin/bash

# clear the output files
out1="PipelineOut1.txt" # small db A
out2="PipelineOut2.txt" # large db A (to ensure difference is not relative to database size)
out3="PipelineOut3.txt" # small db AAAA
out4="PipelineOut4.txt" # small db MX
out5="PipelineOut5.txt" # small db NS
out6="PipelineOut6.txt" # small db TXT

echo '' > "$out1"
echo '' > "$out2"
echo '' > "$out3"
echo '' > "$out4"
echo '' > "$out5"
echo '' > "$out6"




for i in {1..30} ;
do

echo "========== RUNNING ========" >> "$out1"
# start a new name-server
/usr/bin/time -v ./../../main/main 10 10 [::]:10000 4GB >>"$out1" 2>&1 &
sleep 3 
/usr/bin/time -f "Elapsed: %E seconds" bash -c '
dig @127.0.0.1 -p 10000 "20231001.google.nu.history.openintel.nl" A

while true; do
    sleep 0.010
    output=$(dig @127.0.0.1 -p 10000 "20231001.google.nu.history.openintel.nl" A)

    if ! grep -q "HINFO \"WAIT\"" <<< "$output"; then
        break
    fi
done
' >> "$out1" 2>&1

# kill the name server
pkill main  

done


for i in {1..30} ;
do

echo "========== RUNNING ========" >> "$out2"
# start a new name-server
/usr/bin/time -v ./../../main/main 10 10 [::]:10000 4GB >>"$out2" 2>&1 &
sleep 3 
/usr/bin/time -f "Elapsed: %E seconds" bash -c '
dig @127.0.0.1 -p 10000 "20231001.google.fr.history.openintel.nl" A

while true; do
    sleep 0.010
    output=$(dig @127.0.0.1 -p 10000 "20231001.google.fr.history.openintel.nl" A)

    if ! grep -q "HINFO \"WAIT\"" <<< "$output"; then
        break
    fi
done
' >> "$out2" 2>&1

# kill the name server
pkill main  

done


for i in {1..30} ;
do

echo "========== RUNNING ========" >> "$out3"
# start a new name-server
/usr/bin/time -v ./../../main/main 10 10 [::]:10000 4GB >>"$out3" 2>&1 &
sleep 3 
/usr/bin/time -f "Elapsed: %E seconds" bash -c '
dig @127.0.0.1 -p 10000 "20231001.google.nu.history.openintel.nl" AAAA

while true; do
    sleep 0.010
    output=$(dig @127.0.0.1 -p 10000 "20231001.google.nu.history.openintel.nl" AAAA)

    if ! grep -q "HINFO \"WAIT\"" <<< "$output"; then
        break
    fi
done
' >> "$out3" 2>&1

# kill the name server
pkill main  

done


for i in {1..30} ;
do

echo "========== RUNNING ========" >> "$out4"
# start a new name-server
/usr/bin/time -v ./../../main/main 10 10 [::]:10000 4GB >>"$out4" 2>&1 &
sleep 3 
/usr/bin/time -f "Elapsed: %E seconds" bash -c '
dig @127.0.0.1 -p 10000 "20231001.google.nu.history.openintel.nl" MX

while true; do
    sleep 0.010
    output=$(dig @127.0.0.1 -p 10000 "20231001.google.nu.history.openintel.nl" MX)

    if ! grep -q "HINFO \"WAIT\"" <<< "$output"; then
        break
    fi
done
' >> "$out4" 2>&1

# kill the name server
pkill main  

done


for i in {1..30} ;
do

echo "========== RUNNING ========" >> "$out5"
# start a new name-server
/usr/bin/time -v ./../../main/main 10 10 [::]:10000 4GB >>"$out5" 2>&1 &
sleep 3 
/usr/bin/time -f "Elapsed: %E seconds" bash -c '
dig @127.0.0.1 -p 10000 "20231001.google.nu.history.openintel.nl" NS

while true; do
    sleep 0.010
    output=$(dig @127.0.0.1 -p 10000 "20231001.google.nu.history.openintel.nl" NS)

    if ! grep -q "HINFO \"WAIT\"" <<< "$output"; then
        break
    fi
done
' >> "$out5" 2>&1

# kill the name server
pkill main  

done


for i in {1..30} ;
do

echo "========== RUNNING ========" >> "$out6"
# start a new name-server
/usr/bin/time -v ./../../main/main 10 10 [::]:10000 4GB >>"$out6" 2>&1 &
sleep 3 
/usr/bin/time -f "Elapsed: %E seconds" bash -c '
dig @127.0.0.1 -p 10000 "20231001.google.nu.history.openintel.nl" TXT

while true; do
    sleep 0.010
    output=$(dig @127.0.0.1 -p 10000 "20231001.google.nu.history.openintel.nl" TXT)

    if ! grep -q "HINFO \"WAIT\"" <<< "$output"; then
        break
    fi
done
' >> "$out6" 2>&1

# kill the name server
pkill main  

done