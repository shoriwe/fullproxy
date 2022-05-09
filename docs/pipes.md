# Pipes

Pipes are used to separate the logic of proxying protocols and network connections.

## Bind pipe

### Successful negotiation

```mermaid
sequenceDiagram
Client ->> Proxy : Connects
Proxy -> Client : Protocol negotiation
Proxy ->> Server : Stablish connection
Client --> Server : Proxy traffic
```

## Master / Slave pipe

This pipe is intended to forward traffic accessible by a machine that cannot bind but can connect to external servers,
for this same reason this kind of proxy pipe doesn't support commands like **SOCKS5 BIND**.

```mermaid
sequenceDiagram
participant Client
participant Master
participant Slave
Participant Server
Slave ->> Master : Stablish command connection
Client ->> Master : Connects
Client -> Master : Protocol negotiation
Master ->> Slave : Sends connect command
Slave -->> Master : Stablish new forward traffic connection
Slave ->> Server : Stablish connection
Client --> Server : Proxy traffic
```