def basic_login(username, password)
    if username == "sulcud" and password == "password"
        return true
    end
    return false
end

set_auth(basic_login)