#!/bin/bash

# clear the output files
out1="DuckDBQueryOut1.txt" # small db A
out2="DuckDBQueryOut2.txt" # large db A (to ensure difference is not relative to database size)
out3="DuckDBQueryOut3.txt" # small db AAAA
out4="DuckDBQueryOut4.txt" # small db MX
out5="DuckDBQueryOut5.txt" # small db NS
out6="DuckDBQueryOut6.txt" # small db TXT

echo '' > "$out1"
echo '' > "$out2"
echo '' > "$out3"
echo '' > "$out4"
echo '' > "$out5"
echo '' > "$out6"


for i in {1..30};
do
    /usr/bin/time -v ./duckdb_query 1 
    sleep 1
done >> "$out1" 2>&1

for i in {1..30};
do
    /usr/bin/time -v ./duckdb_query 2
    sleep 1
done >> "$out2" 2>&1

for i in {1..30};
do
    /usr/bin/time -v ./duckdb_query 3
    sleep 1
done >> "$out3" 2>&1

for i in {1..30};
do
    /usr/bin/time -v ./duckdb_query 4
    sleep 1
done >> "$out4" 2>&1

for i in {1..30};
do
    /usr/bin/time -v ./duckdb_query 5
    sleep 1
done >> "$out5" 2>&1

for i in {1..30};
do
    /usr/bin/time -v ./duckdb_query 6
    sleep 1
done >> "$out6" 2>&1
