def basic_login(username, password)
    if username == "sulcud" and password == "password"
        return True
    end
    return False
end

SetAuth(basic_login)