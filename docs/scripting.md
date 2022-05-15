# Scripting

This document describe the interfaces used in `fullproxy` driver development.

## Authentication functions

Functions intended for authentication will be loaded using `SetAuth(function)`, functions should expect two string arguments, specifically corresponding to username and password, authentication succeed will be evaluated by the return values `True` in case of success or `False` in case of failure, notice if the function call raises an error it will be also considered an auth failure.

Example:

```ruby
def basic_login(username, password)
    if username == "sulcud" and password == "password"
        return True
    end
    return False
end

SetAuth(basic_login)
```

## Filtering functions

Function intended to filter incoming, outgoing, listens and accepts can be loaded in drivers using:

- `SetInbound(function)`
- `SetOutbound(function)`
- `SetListen(function)`
- `SetAccept(function)`

This functions are expected to receive a string value containing the `HOST:PORT` value of the connection. Allowing the connection will be evaluated by the return values `True` in case of success or `False` in case of failure, notice if the function call raises an error it will be also considered an allow failure.

Example:

```ruby
def no_localhost(address)
    if "127.0.0.1" in address
        return False
    elif "localhost" in address
        return false
    end
    return True
end

SetInbound(no_localhost)
SetOutbound(no_localhost)
SetListen(no_localhost)
SetAccept(no_localhost)
```

