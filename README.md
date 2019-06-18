# rdma-ds

This dameon-set will run on every node and will be a RESTful API endpoint for querying data about the node that it is currently running on.

## Test API
To test the API run the following commands:
```
cd v1
go test
```

## Client
To use the client library import `github.com/swrap/rdma-ds/v1` and utilize the built in supported functions.

## Server
To start the server run the script below:
```
./build_run_server.sh
```
If you want to specify a port number, set the environment variable `PORT` like below:
```
PORT=40007 ./build_run_server.sh
```
