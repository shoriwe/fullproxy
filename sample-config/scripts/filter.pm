def no_localhost(address)
    if "127.0.0.1" in address
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