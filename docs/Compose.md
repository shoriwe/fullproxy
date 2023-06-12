# Compose

This document describes the `compose` contract specification.

## Root structure

### Dependencies

- [Circuits](#Circuits)
- [Proxies](#Proxies)
- [Slaves](#Slaves)

### Definition

```yaml
circuits:
    NAME:
        <<CIRCUIT>>
proxies:
    NAME:
        <<PROXY>>
slaves:
    NAME:
        <<SLAVE>>
```

### Examples

```yaml
circuits:
    australia:
        <<CIRCUIT>>
    china:
        <<CIRCUIT>>
proxies:
    US:
        <<PROXY>>
slaves:
    colombia:
        <<SLAVE>>
    us:
        <<SLAVE>>
```

## Network

### Dependencies

- [Network](#Network)
- [Auth](#Auth)
- [Crypto](#Crypto)
- [Filter](#Filter)

### Definition

```yaml
type: basic|master|ssh
network: tcp|EMPTY
address: HOST:PORT|EMPTY
data:
    <<NETWORK>>|EMPTY
control:
    <<NETWORK>>|EMPTY
auth:
    <<AUTH>>|EMPTY
crypto:
    <<CRYPTO>>|EMPTY
slaveListener: true|false
listenerFilter:
    <<FILTER>>|EMPTY
dialFilter:
    <<FILTER>>|EMPTY
```

### Examples

- Basic listener

```yaml
type: basic
network: tcp
address: 0.0.0.0:443
```

- Master listener

```yaml
type: master
data:
    type: basic
    network: tcp
    address: 0.0.0.0:9050
control:
    type: basic
    network: tcp
    address: 10.10.50.10:8000
```

- SSH listener

```yaml
type: ssh
network: tcp
address: 10.10.50.10:22
auth:
    username: sulcud
    password: password
```

- Basic listener with **TLS**

```yaml
type: basic
type: basic
network: tcp
address: 0.0.0.0:443
crypto:
    mode: tls
    selfSigned: true
```

- Master with **Data** coming from **Basic** and Slave from TLS **SSH**

```yaml
type: master
data:
    type: basic
    network: tcp
    address: 0.0.0.0:8000
control:
    type: ssh
    network: tcp
    address: 10.10.50.10:22
    auth:
        privateKey: /home/admin/id_rsa
    crypto:
        mode: tls
        selfSigned: true
        insecureSkipVerify: true
```

## Auth

### Definition

```yaml
username: "STRING"|EMPTY
password: "STRING"|EMPTY
privateKey: "STRING"|EMPTY
serverKey: "STRING"|EMPTY
```

### Examples

- Username and Password

```yaml
username: sulcud
password: password
```

- Private Key

```yaml
privateKey: /path/to/private/key
```

- Server key

```yaml
serverKey: /path/to/server/key
```

## Crypto

### Definition

```yaml
mode: tls
selfSigned: true|false
insecureSkipVerify: true|false
cert: /path/to/crt
key: /path/to/key
```

### Examples

- Self signed key

```yaml
mode: tls
selfSigned: true
```

- Insecure skip verify

```yaml
mode: tls
insecureSkipVerify: true
```

- Cert and key from file

```yaml
mode: tls
cert: /etc/certs/crt
key: /etc/certs/key
```

## Circuits

### Dependencies

- [Network](#Network)
- [Knots](#Knots)

### Definition

```yaml
network: tcp
address: HOST:PORT
listener:
    <<NETWORK>>
knots:
    - <<KNOT>>
```

### Examples

```yaml
network: tcp
address: google.com:443
listener:
    type: basic
    network: tcp
    address: localhost:443
knots:
    - type: ssh
      network: tcp
      address: 10.10.50.10:22
      auth:
          username: sulcud
          password: password
    - type: socks5
      network: tcp
      address: 207.208.10.30:9050
```

## Knots

### Dependencies

- [Auth](#Auth)

### Definition

```yaml
type: forward|socks5|ssh
network: tcp
address: HOST:PORT
auth:
    <<AUTH>>|EMPTY
```

### Examples

```yaml
type: ssh
network: tcp
address: google.com:22
auth:
    username: sulcud
    password: password
```

## Filter

### Dependencies

- [Match](#Match)

### Definition

```yaml
whitelist:
    - <<MATCH>>
blacklist:
    - <<MATCH>>
```

### Examples

```yaml
whitelist:
    - <<MATCH>>
blacklist:
    - <<MATCH>>
```

## Match

### Definition

```yaml
host: REGEXP|EMPTY
port: NUMBER|EMPTY
portRange:
    from: NUMBER|EMPTY
    to: NUMBER|EMPTY
```

### Examples

```yaml
host: "127\\.0\\.0\\.\\d+"
port: 443
portRange:
    from: 0
    to: 65535
```

## Proxies

### Dependencies

- [Network](#Network)
- [Auth method](#Auth method)

### Definition

```yaml
type: forward|http|socks5
listener:
    <<NETWORK>>
dialer:
    <<NETWORK>>|EMPTY
network: tcp|EMPTY
address: HOST:PORT|EMPTY
authMethods:
    - <<AUTH_METHOD>>
```

### Examples

- Forward

```yaml
type: forward
listener:
    type: basic
    network: tcp
    address: 0.0.0.0:80
network: tcp
address: google.com:80
```

- Socks5 with dialer in SSH server

```yaml
type: socks5
listener:
    type: basic
    network: tcp
    address: 0.0.0.0:9050
dialer:
    type: ssh
    network: tcp
    address: 10.10.50.10:22
    auth:
        privateKey: /home/admin/id_rsa
```

### Auth method

### Definition

```yaml
raw:
    USERNAME: PASSWORD
```

### Examples

```yaml
raw:
    sulcud: password
    shoriwe: password
```

## Slaves

### Dependencies

- [Network](#Network)

### Definition

```yaml
masterNetwork: tcp
masterAddress: HOST:PORT
dialer:
    <<NETWORK>>
listener:
    <<NETWORK>>|EMPTY
```

### Examples

```yaml
masterNetwork: tcp
masterAddress: 10.10.50.10:9999
dialer:
    <<NETWORK>>
listener:
    <<NETWORK>>
```