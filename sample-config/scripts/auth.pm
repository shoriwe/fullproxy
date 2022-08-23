def basic_login(username, password)
    if username == "sulcud" and password == "password"
        return True
    end
    return False
end

set_auth(basic_login)