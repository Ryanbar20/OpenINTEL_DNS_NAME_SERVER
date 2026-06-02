#!/bin/bash

# start the name-server
./../main/main 10 10 [::]:10000 >/dev/null &
pid=$!
sleep 3 


#clear output file
echo '' > out2.txt 
echo "============================================================" >> out2.txt
echo 'submit some invalid queries. Observe that they are refused' >> out2.txt
echo "============================================================" >> out2.txt
dig @127.0.0.1 -p 10000 "20201001.google.nu.history.openintel." A >> out2.txt
dig @127.0.0.1 -p 10000 "20201001.google.nu.histry.openintel.nl" AAAA >> out2.txt
dig @127.0.0.1 -p 10000 "202010ABC01.google.nu.history.openintel.nl" NS >> out2.txt
dig @127.0.0.1 -p 10000 "20201001.history.openintel.nl" MX >> out2.txt
dig @127.0.0.1 -p 10000 "20201001.google.nu.history.closedintel.nl" TXT >> out2.txt

# kill the name server
kill $pid 