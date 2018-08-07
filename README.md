# sbps

sbps (Simple Broadcast Proxy Server) is the broadcast proxy server between servers and clients.

![sbps architecture](img/sbps_architecture.PNG)

A server resource represents a server or a resource that sbps could read from or write to. A client connection represents a connection between sbps and clients, receives data from servers and sends data to servers. sbps makes dedicated goroutines for each server resources and clients connections. Each goroutines monitors a server resource or a client connection state. If a server resource or a client connection is closed, a related goroutine stopped. Server resource goroutines and client connection goroutines communicate with each other directly. Server resource and Client connection goroutines are in N:M relationship.

## Option

## Usage Examples
* TCP
~~~
# sbps -mode TCP:5000 -resource TCP:192.168.0.200:5000:RW -interval 4
~~~

* UDP with write mode, Unix with read mode, No retry
~~~
# sbps -mode TCP:6000 -resource UDP:192.168.0.200:5000:W,UNIX:/root/sbps_res:R -interval 0
~~~

* FIFO with read mode, FIFO with write mode
~~~
# sbps -mode UNIX:/root/sbps_server -resource FIFO:/root/sbps_fifo_r:R,FIFO:/root/sbps_fifo_w:W -interval 2
~~~

* Broadcast echo server
~~~
# sbps -mode UNIX:/root/sbps_server -resource FIFO:/root/sbps_fifo:RW -interval 3
~~~

## Build and run

* Set Env
~~~
# export GOPATH=...
# export GOBIN=...
~~~

* Get and build
~~~
# go get github.com/ssup2/sbps/
# cd $GOPATH/src/github.com/ssup2/sbps/
# make
~~~

* Run
~~~
# cd GOBIN
# ./sbps ...
~~~


