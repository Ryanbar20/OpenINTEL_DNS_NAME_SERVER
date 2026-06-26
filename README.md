# OpenINTEL DNS name-server
This repository presents a Golang module for a DNS name-server implementation for querying the OpenINTEL dataset.

## Usage
To use the module, simply add `import "github.com/Ryanbar20/OpenINTEL_DNS_NAME_SERVER"` into your go file.
For a template of how to do this, see the `main/main.go` file.

## Navigation
The root of the repository contains the module for the DNS name-server.
The `main` folder contains a sample module for running the name-server.
The `systemTest` folder contains Bash files for running system tests for the name-server. Please refer to `systemtest/README.md` for more information.
The `benchmarks` folder contains directories for some quantitative tests against which the name-server can be tested. Please refer to `benchmarks/README.md` for more information.

## Version
The newest version of the code (which contains the root and `main` folders) is presented in version 0.0.23 (see tag v0.0.23). All commits after this solely added documentations and the tests in the `benchmarks` and `systemTest` folders.