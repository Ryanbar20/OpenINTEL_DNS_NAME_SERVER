#!/bin/bash

# start the name-server
./../main/main 10 10 [::]:10000 >/dev/null &
pid=$!
sleep 3 

#clear output file
echo '' > out1.txt 

echo "============================================================" >> out1.txt
echo 'submit an A query and see a wait/limit in return' >> out1.txt
echo "============================================================" >> out1.txt
dig @127.0.0.1 -p 10000 "20201001.google.nu.history.openintel.nl" A >> out1.txt

echo "============================================================" >> out1.txt
echo 'submit an AAAA query and see a wait/limit in return' >> out1.txt
echo "============================================================" >> out1.txt
dig @127.0.0.1 -p 10000 "20201001.google.nu.history.openintel.nl" AAAA >> out1.txt

echo "============================================================" >> out1.txt
echo 'submit an NS query and see a wait/limit in return' >> out1.txt
echo "============================================================" >> out1.txt
dig @127.0.0.1 -p 10000 "20201001.google.nu.history.openintel.nl" NS >> out1.txt

echo "============================================================" >> out1.txt
echo 'submit an MX query and see a wait/limit in return' >> out1.txt
echo "============================================================" >> out1.txt
dig @127.0.0.1 -p 10000 "www.20201001.google.nu.history.openintel.nl" MX >> out1.txt

echo "============================================================" >> out1.txt
echo 'submit a TXT query and see a wait/limit in return' >> out1.txt
echo "============================================================" >> out1.txt
dig @127.0.0.1 -p 10000 "20201001.google.nu.history.openintel.nl" TXT >> out1.txt


# kill the name server
kill $pid
