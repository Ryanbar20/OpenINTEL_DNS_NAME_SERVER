#!/bin/bash

# start the name-server
./../main/main 10 10 [::]:10000 &
pid=$!
sleep 3 
# submit some invalid queries
dig @127.0.0.1 -p 10000 "20201001.google.nu.history.openintel." A > out2.txt
dig @127.0.0.1 -p 10000 "20201001.google.nu.histry.openintel.nl" AAAA >> out2.txt
dig @127.0.0.1 -p 10000 "202010ABC01.google.nu.history.openintel.nl" NS >> out2.txt
dig @127.0.0.1 -p 10000 "20201001.history.openintel.nl" MX >> out2.txt
dig @127.0.0.1 -p 10000 "20201001.google.nu.history.closedintel.nl" TXT >> out2.txt

# kill the name server
kill $pid 
