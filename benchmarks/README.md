# Benchmarks
This folder contains relevant benchmarks / quantitative tests for validation of the name-server implementation.

## database
This folder contains a test for the speed and memory usage of querying the OpenINTEL dataset via the name-server, compared to a direct DuckDB implementation.
This test validates that the name-server uses a dataset querying implementation alike the implementation given in `duckdb_query.go`.

## multiThread
This folder contains a test for the Queries per Second that a regular pipeline achieves, compared to that of a pipeline that is busy processing a database query.
This test validates that the multi threaded approach is correctly implemented, and that the database-querying does not have a great impact on the performance of the cache-hit part of the pipeline.

## name_server
This folder contains a simple MiekgDNS implementation of a name server. This implementation can only answer 3 queries, which are the same as the cache-hit queries in the `multiThread` tests.
This test measures the Queries per Second that the implementatin achieves. It allows comparing the pipelines in the `multiThread` tests with a regular MiekgDNS implementation, to see how well the cache-hit part of the pipeline performs.


## Usage
To run the tests, make sure that the `benchmarks/database/duckdb_query.go`, `benchmarks/name_server/miekg.go` and `main/main.go` files are compiled to corresponding executables (that have the same path, but the .go extension is left out).
Next to this, `dnsperf` has to be installed to run the `multiThread` and `name_server` tests.