#!/bin/bash

if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root"
   echo "This is because it needs to create an additional IP loopback address to locally send dns queries from"
   echo "This address will be deleted afterwards"
   exit 1
fi

bash ./test1.sh
echo 'test 1 done'

bash ./test2.sh
echo 'test 2 done'

bash ./test3.sh
echo 'test 3 done'

bash ./test4.sh
echo 'test 4 done'

bash ./test5.sh
echo 'test 5 done'

bash ./test6.sh
echo 'test 6 done'

bash ./test7.sh
echo 'test 7 done'