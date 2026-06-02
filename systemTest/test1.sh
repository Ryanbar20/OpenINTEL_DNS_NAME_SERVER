#!/bin/bash

# start the name-server
./../main/main 10 10 [::]:10000 &
pid=$!
sleep 3 
# submit an A query and save output
dig @127.0.0.1 -p 10000 "20201001.google.nu.history.openintel.nl" A > out1.txt
# submit an AAAA query and save output
dig @127.0.0.1 -p 10000 "20201001.google.nu.history.openintel.nl" AAAA >> out1.txt
# submit an NS query and save output
dig @127.0.0.1 -p 10000 "20201001.google.nu.history.openintel.nl" NS >> out1.txt
# submit an MX query and save output
dig @127.0.0.1 -p 10000 "www.20201001.google.nu.history.openintel.nl" MX >> out1.txt
# submit a TXT query and save output
dig @127.0.0.1 -p 10000 "20201001.google.nu.history.openintel.nl" TXT >> out1.txt


# kill the name server
kill $pid 
