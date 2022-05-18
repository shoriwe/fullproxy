# YAML config

```yaml
init-order:
  - LISTENER_NAME
  - ...
drivers: # This field is optional.
  DRIVER_NAME: /PATH/TO/PLASMA/SCRIPT
listeners:
  LISTENER_NAME:
    log: /PATH/TO/FILE/TO/DATA
    sniff:
      incoming: /PATH/TO/FILE/WITH/INCOMING/TRAFFIC
      outgoing: /PATH/TO/FILE/WITH/OUTGOING/TRAFFIC
    config:
      # Mandatory by all types of listeners
      type: basic | master | slave
      network: tcp | unix
      address: HOST:PORT | /PATH/TO/UNIX/SOCK
      tls: # Ignore to generate a self signed cert
        - /PATH/TO/TLS/PEM:/PATH/TO/TLS/KEY
        - ...
      # Available for all types of protocol
      filters:
        inbound: DRIVER_NAME # Ignore to no filer
        outbound: DRIVER_NAME # Ignore to no filer
        listen: DRIVER_NAME # Ignore to no filer
        accept: DRIVER_NAME # Ignore to no filer

      # Mandatory by master and slave
      master-network: tcp | unix
      master-address: HOST:PORT | /PATH/TO/UNIX/SOCK

      # Mandatory by master
      master-tls: # Ignore to generate a self signed cert.
        - /PATH/TO/TLS/PEM:/PATH/TO/TLS/KEY

      # Mandatory by slave
      slave-trust: true | false

    protocol: # Used only when type is basic | master
      # Mandatory
      type: socks5|http|reverse-raw|reverse-http|forward|translate|http-hosts

      # Only for socks5 and http
      authentication: DRIVE_NAME # Ignore to no auth

      # Mandatory by forward and translate
      target-network: tcp | unix
      target-address: HOST:PORT | /PATH/TO/UNIX/SOCK

      # Mandatory by translate
      proxy-network: tcp | unix
      proxy-address: HOST:PORT | /PATH/TO/UNIX/SOCK
      translation: socks5:forward # Currently only supported
      credentials: USERNAME:PASSWORD

      # Mandatory for reverse-raw
      raw-hosts:
        NAME:
          network: tcp | unix
          address: HOST:PORT | /PATH/TO/UNIX/SOCK

      # Mandatory for reverse-http
      http-hosts:
        HOSTNAME:
          path: /wanted/uri
          response-headers: # Headers to in inject in the response to the client
            - KEY:VALUE
          request-headers: # Headers to in inject in the request to the server
            - KEY:VALUE
          pool: # Load balancing pool
            NAME:
              url: URL
              network: tcp | unix
              address: HOST:PORT | /PATH/TO/UNIX/SOCK
```
