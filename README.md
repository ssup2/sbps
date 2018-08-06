# sbps

sbps (Simple Broadcast Proxy Server) is the broadcast proxy server between servers and clients.

![sbps architecture](img/sbps_architecture.PNG)

A server resource represents a server or a resource that sbps could read from or write to. A client connection represents a connection between sbps and clients, receives data from servers and sends data to servers. sbps makes dedicated goroutines for each server resources and clients connections. Each goroutines monitors a server resource or a client connection state. If a server resource or a client connection is closed, a related goroutine stopped. Server resource goroutines and client connection goroutines communicate with each other directly. Server resource and Client connection goroutines are in N:M relationship.
 
