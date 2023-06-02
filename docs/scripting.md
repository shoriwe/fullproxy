# Scripting

This document describe the interfaces used in `fullproxy` driver development.

## Authentication functions

Functions intended for authentication will be loaded using `set_auth(function)`, functions should expect two string
arguments, specifically corresponding to username and password, authentication succeed will be evaluated by the return
values `True` in case of success or `False` in case of failure, notice if the function call raises an error it will be
also considered an auth failure.

Example:

```ruby
def basic_login(username, password)
    if username == "sulcud" and password == "password"
        return True
    end
    return False
end

set_auth(basic_login)
```

## Filtering functions

Function intended to filter incoming, outgoing, listens and accepts can be loaded in drivers using:

- `set_inbound(function)`
- `set_outbound(function)`
- `set_listen(function)`
- `set_accept(function)`

This functions are expected to receive a string value containing the `HOST:PORT` value of the connection. Allowing the
connection will be evaluated by the return values `True` in case of success or `False` in case of failure, notice if the
function call raises an error it will be also considered an allow failure.

Example:

```ruby
def no_localhost(address)
    if "localhost" in address
        return False
    elif "localhost" in address
        return false
    end
    return True
end

set_inbound(no_localhost)
set_outbound(no_localhost)
set_listen(no_localhost)
set_accept(no_localhost)
```
