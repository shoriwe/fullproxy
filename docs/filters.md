# Filters

## Incoming connections

```
sequenceDiagram
Client ->> Proxy : Connects
alt Client address whitelisted
Proxy -> Client : Protocol negotiation
else Client address blacklisted
Proxy -->> Client : Closes connection
end
```

## Outgoing connections

```
sequenceDiagram
Client ->> Proxy : Connects
Client -> Proxy : Protocol Negotiation
alt Server address whitelisted
Proxy ->> Server : Stablish connection
Client --> Server : Proxy traffic
else Server address blacklisted
Proxy -->> Client : Closes connection
end
```

## Listen

```
sequenceDiagram
Client ->> Proxy : Connects
Client -> Proxy : Protocol negotiation
alt Listen address whitelisted
Proxy -> Proxy : Bind address
else Listen address blacklisted
Proxy -->> Client : Closes connection
end
```

## Accept

```
sequenceDiagram
Client ->> Proxy : Connects
Client -> Proxy : Protocol negotiation
Proxy ->> Proxy : Bind address
Client 2 ->> Proxy : Connects
alt Client 2 address whitelisted
Client --> Client 2 : Proxy traffic
else Client 2 address blacklisted
Proxy -->> Client 2 : Closes connection
Proxy -->> Client : Closes connection
end
```