# System tests
This directory contains bash scripts for executing relevant system tests for the name-server.
The provided system tests test the following requirements identified in the corresponding paper. 
The results of each script may depend on network speed, as data downloads are expected to be done in a certain amount of time.
When there are unexpected results for this reason, the amount of time that the script waits should be increased.


- [ ] TODO#1 Refer to the corresponding paper.


## Test 1
### Description
In test 1, DNS queries for the supported query types (A, AAAA, TXT, NS and MX) should be submitted to the name server. After each query, some time should be waited, after which a retry can be sent. This retry, or later retries, will eventually be answered by the answers.
### Test target
This test asserts that the name server supports these query types and sends correctly formatted DNS messages back. It also asserts that the server sends `wait` messages when it should
### Expected results
The results of this test are that each query gets answered by a `wait` message, as the answers to these queries are not in the cache. The retries will return the actual answer in a correctly formatted DNS message.

## Test 2
### Description
In Test 2, DNS queries are to be sent to the name-server. These queries should contain a question with an invalid format or an usupported type.
### Test target
Test 2 asserts that the name server correctly handles the question format `YYYYMMDD.<domain>.history.openintel.nl`. It also asserts that the server refuses unsupported types.
### Expected results
The results of this test should be that all queries with an invalidly formatted question or unsupported are refused by the name server.

## Test 3
### Description
In Test 3, an A query for google.nu on October 1st 2023 is to be sent to the name-server, after which a `wait` message is sent. After the first retry that returns an answer, all following retries will return this same answer.
### Test target
Test 3 asserts that the question in the DNS query is correctly parsed and that it queries the OpenINTEL dataset with the parameters from the query. Test 3 also asserts that the name-server correctly stores an answer to a question in its cache and uses this cache for answerring the retries.
### Expected results
The first question is answerred with a `wait` message. The first retry that returns an answer will indicate ip `172.217.23.196`, any following retries will give the same answer.

## Test 4
### Description
In Test 4, a server with cache size 1 is started. A query is to be sent to this name-server and there should be waited until retries for this query return the answer to the query. This same process is then done with a different query, after which the first query is retried.
### Test target
This test asserts that the server correctly implements the cache by checking that it maintains the cache size limit.
### Expected results
The final retry for the first query, after waiting for the second query, should be answerred by a `wait` message again. This is because the name server has a cache size of 1, and thus the oldest entry is dropped from the cache, after a new query is answered.

## Test 5
### Description
In Test 5, a server is started with queue size 2. 3 distinct queries from distinct source IPs are to be sent to the server. Then, send a different cache-miss query from the first IP.
### Test target
This test asserts that the server correctly implements the queue and wait/limit messages.
### Expected results
The first 2 queries are answered with a `wait` message. The 3rd query is answered with a `limit` message that indicates that the queue is full. The 4th query is also answered with a `limit` message, but this time it indicates that the user already has a query in the queue.

## Test 6
### Description
In test 6, 2 identical cache-miss queries are to be sent to the name-server from distinct IPs. Then, a different query is to be sent from the second IP.
### Test target
This test asserts that the server correctly implements the queue and wait/limit messages.
### Expected results
The first 2 queries are answered with a `wait` message, that both indicate the same queue position, as these two queries now 'share' that position. The third query is answered with a `wait` message too, but this time it has a higher queue position. This is because the first IP sent the first query, thus the second IP did not sent a new query for the queue itself. For this reason, the second IP can still send a new query for the queue.

## Test 7
### Description
In test 7, a query is to be sent to this name-server and there should be waited until retries for this query return the answer to the query. Then, Another query is sent and immediately the first query is re-sent. At the end, the second query is retried until an answer is sent back.
### Test target
This test asserts that the server correctly implements a multithreaded design with separation between cache-miss and cache-hit queries.
### Expected results
The first query will be answered. After this, the second query is sent and the first query is re-sent. This retransmission of the first query gets answered befor the second query is done processing. This is asserted by seeing that retries of the second query still return `wait` messages after the retransmission of the first query is answered.