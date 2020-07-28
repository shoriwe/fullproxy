# FullProxy
 Bind and reverse connection based, SOCKS5 and HTTP proxy
 
 # Description
 ![FullProxyLogo](https://raw.githubusercontent.com/shoriwe/FullProxy/master/logo/full-proxy-logo.PNG) \
`FullProxy` is a `Bind` and `Reverse Connection` based `HTTP` and `SOCKS5` portable proxy
# Usage
## Implemented protocols
```
user@linux:~$ fullproxy help
Usage: FullProxy PROTOCOL *FLAGS
Protocols available:
         - socks5
         - http
         - interface-master
```
### SOCKS5 protocol
```
user@linux:~$ fullproxy socks5 --help
Usage of socks5:
  -interface-mode
        Connect to an interface, no bind proxying
  -ip string
        IP address to listen on. When "-interface-mode" flag is set, is the IP of interface to connect (default "0.0.0.0")
  -password string
        Password of the running proxy, requires "-username" and can't be an empty string ('')
  -port string
        Port address to listen on. When "-interface-mode" flag is set, is the Port of the interface to connect (default "1080")
  -username string
        Username of the running proxy, requires "-password" and can't be an empty string ('')
```
### HTTP protocol
NOT IMPLEMENTED YET
### Interface master protocol
```
Usage of interface-master:
  -ip string
        IP address to listen on. (default "0.0.0.0")
  -port string
        Port address to listen on. (default "1080")
```
# Concepts
## Interface mode
Handles the proxying between a reverse connected proxy and the clients. In other words, it will receive the connections of the clients and will forward the traffic to the proxy that is reverse connected to it.
### How it works
1. It first binds to the address specified by the user.
2. Then accept the connection from the proxy server.
3. Finally, it proxy the traffic of all new incoming connections to the proxy server that was reverse connected to it in the second step.
In other words, is the proxy of another proxy but totally invisible for the client.
### Applications
This could be specially useful when you need to proxy a network that a machine have access to, but you can't bind with it
### Considerations
- The `interface` protocol may loss some setup connections if it is extremely stressed, but it should work `just fine` if the connections where already made
