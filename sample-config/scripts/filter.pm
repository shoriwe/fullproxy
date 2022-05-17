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