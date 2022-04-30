# socks5

This is a Golang implementation of the Socks5 protocol library.   
To see in this [SOCKS Protocol Version 5](https://www.rfc-editor.org/rfc/rfc1928.html).  
This library is also compatible with Socks4 and Socks4a.

# Contents

- [Features](#Features)
- [Install](#Installation)
- [Examples](#Examples)
    - [Server example](#Server-example)
        - [simple (no authentication)](#simple-no-authentication)
        - [username/password authentication in memory](#username/password-authentication-in-memory)
        - [custom transporter to transmit data between client and remote](#custom-transporter-to-transmit-data-between-client-and-remote)
    - [Client example](#Client)
        - [CONNECT usage](#CONNECT-usage)
        - [UDP_ASSOCIATE usage](#UDP_ASSOCIATE-usage)
        - [BIND usage](#BIND-usage)
- [FAQ](#FAQ)

# Features

- socks5:
    - command: **CONNECT**, **UDP ASSOCIATE**, **BIND**.
    - auth methods:
        - **Username/Password** authentication.
        - No Authentication Required.
- socks4:
    - command: **CONNECT**, **BIND**.
    - auth: (no support).
- sock4a: same as socks4.
- Custom client and server authenticator.
- Easy to read source code.
- Similar to the Golang standard library experience.

# Installation

``` sh
$ go get "github.com/haochen233/socks5"`
```

# Examples

## Server example

### simple (no authentication):

```go
package main

import (
	"log"
	"github.com/haochen233/socks5"
)

func main() {
	// create socks server.
	srv := &socks5.Server{
		// socks server listen address.
		Addr: "127.0.0.1:1080",
		// UDP assocaite and bind command listen ip.
		// Don't need port, the port will automatically chosen.
		BindIP: "127.0.0.1",
		// if nil server will provide no authentication required method.
		Authenticators: nil,
	}

	// start listen
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}


```

### username/password authentication in memory:

```go
package main

import (
	"crypto/md5"
	"log"

	"github.com/haochen233/socks5"
)

func main() {
	// create a username/password store in memory.
	var userStorage socks5.UserPwdStore = socks5.NewMemeryStore(md5.New(), "secret")
	// set a pair of username/password.
	userStorage.Set("admin", "123456")

	srv := &socks5.Server{
		Addr:   "127.0.0.1:1080",
		BindIP: "127.0.0.1",
		// enable username/password method and authenticator.
		Authenticators: map[socks5.METHOD]socks5.Authenticator{
			socks5.USERNAME_PASSWORD: socks5.UserPwdAuth{UserPwdStore: userStorage},
			// There is already an authentication method.
			// If want enable no authentication required method.
			// you should enable it explicit.
			socks5.NO_AUTHENTICATION_REQUIRED: socks5.NoAuth{},
		},
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
```

### custom transporter to transmit data between client and remote.

```go
package main

import (
	"log"
	"net"

	"github.com/haochen233/socks5"
)

// simulate to impl socks5.Transporter interface.
// transport encrypted data.
type cryptTransport struct {
}

func (c *cryptTransport) TransportTCP(client *net.TCPConn, remote *net.TCPConn) <-chan error {
	//encrypt data and send to remote
	//decrypt data and send to client
	return nil
}

func (c *cryptTransport) TransportUDP(server *socks5.UDPConn, request *socks5.Request) error {
	panic("implement me")
	return nil
}

func main() {
	server := &socks5.Server{
		Addr:   "127.0.0.1:1080",
		BindIP: "127.0.0.1",
		// replace default Transporter interface
		Transporter: &cryptTransport{},
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
```

## Client example

### CONNECT usage:

```go
package main

import (
	"log"

	"github.com/haochen233/socks5"
)

func main() {
	// create socks client
	clnt := socks5.Client{
		ProxyAddr: "127.0.0.1:1080",
		// Authenticator supported by the client.
		// It must not be nil.
		Auth: map[socks5.METHOD]socks5.Authenticator{
			// If client want send NO_AUTHENTICATION_REQUIRED method to server, must
			// add socks5.NoAuth authenticator explicitly
			socks5.NO_AUTHENTICATION_REQUIRED: &socks5.NoAuth{},
		},
	}

	// client send CONNECT command and get a tcp connection.
	// and use this connection transit data between you and www.google.com:80.
	conn, err := clnt.Connect(socks5.Version5, "www.baidu.com:80")
	if err != nil {
		log.Fatal(err)
	}

	// close connection.
	conn.Close()
}

```

### UDP_ASSOCIATE usage:

```go
package main

import (
	"fmt"
	"log"

	"github.com/haochen233/socks5"
)

func main() {
	clnt := socks5.Client{
		ProxyAddr: "127.0.0.1:1080",
		// client provide USERNAME_PASSWORD method and 
		// NO_AUTHENTICATION_REQUIRED.
		Auth: map[socks5.METHOD]socks5.Authenticator{
			socks5.NO_AUTHENTICATION_REQUIRED: &socks5.NoAuth{},
			socks5.USERNAME_PASSWORD:          &socks5.UserPasswd{Username: "admin", Password: "123456"},
		},
	}

	// client send UDP_ASSOCIATE command and get a udp connection.
	// Empty local addr string a local address (127.0.0.1:port) is automatically chosen.
	// you can specific a address to tell socks server which client address will
	// send udp data. Such as clnt.UDPForward("127.0.0.1:9999").
	conn, err := clnt.UDPForward("")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// send every datagram should add UDP request header.
	someData := []byte("some data")
	// dest addr where are you send to.
	destAddr, _ := socks5.ParseAddress("127.0.0.1:9190")
	// packing socks5 UDP data with dest addr.
	pakcedData, err := socks5.PackUDPData(destAddr, someData)
	// final send you data
	conn.Write(pakcedData)

	// on the contrary.
	// you should unpacked the packet, after received  every packedData.
	buf := make([]byte, 65507)
	conn.Read(buf)

	// unpacking data.
	destAddr, unpackedData, err := socks5.UnpackUDPData(buf)
	// operate your udp data. 
	fmt.Println(unpackedData)
}
```

### BIND usage:

```go
package main

import (
	"encoding/binary"
	"github.com/haochen233/socks5"
	"log"
)

func main() {
	c := socks5.Client{
		ProxyAddr: "172.16.1.28:1080",
		Auth: map[socks5.METHOD]socks5.Authenticator{
			socks5.USERNAME_PASSWORD:          &socks5.UserPasswd{"admin", "123456"},
			socks5.NO_AUTHENTICATION_REQUIRED: &socks5.NoAuth{},
		},
	}

	// connect
	conn1, err := c.Connect(5, "127.0.0.1:9000")
	if err != nil {
		log.Fatal(err)
	}

	dest := "127.0.0.1:9001"
	// bind
	bindAddr, errors, conn2, err := c.Bind(4, dest)
	if err != nil {
		log.Fatal(err)
	}

	// An example tell dest about socks server bind address 
	// via CONNECT proxy connection.
	port := make([]byte, 2)
	binary.BigEndian.PutUint16(port, bindAddr.Port)
	conn1.Write(append(bindAddr.Addr, port...))

	// wait the second reply. if nil the dest already
	// established with socks server.
	err = <-errors
	if err != nil {
		log.Fatal(err)
		return
	}

	// bind success
	_, err = conn2.Write([]byte("hello"))
	if err != nil {
		return
		log.Fatal(err)
	}
}
```

# FAQ:
- Server default enable socks4. How to disable socks4 support?  
  when you initialize a socks5 server, you should spefic this flag to disable explicitly.
   ```go
   server := &socks5.Server{
       DisableSocks4: true,
   }
   ```