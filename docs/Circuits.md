# Circuits

Circuits let user chain proxy protocols to access remote services 

```mermaid
stateDiagram-v2
	direction LR
	
	[*] --> Socks5
	Socks5 --> SSH
	SSH --> [*]
```