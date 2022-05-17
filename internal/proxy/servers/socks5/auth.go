package socks5

func (socks5 *Socks5) UsernamePasswordAuthentication(context *Context) error {
	_ = context.Reply(BasicReply{
		Version:    SocksV5,
		StatusCode: UsernamePassword,
	})
	_, connectionError := context.ClientConnection.Read(context.Chunk[:])
	if connectionError != nil {
		return connectionError
	}
	if context.Chunk[0] != BasicNegotiation {
		return UnsupportedAuthenticationMethod
	}
	userLength := context.Chunk[1]
	username := context.Chunk[2 : 2+userLength]
	passwordLength := context.Chunk[2+userLength]
	password := context.Chunk[2+userLength+1 : 2+userLength+1+passwordLength]
	loginError := socks5.AuthenticationMethod(username, password)
	if loginError != nil {
		_, _ = context.ClientConnection.Write([]byte{BasicNegotiation, FailedAuthentication})
		return loginError
	}
	_, _ = context.ClientConnection.Write([]byte{BasicNegotiation, SucceedAuthentication})
	return nil
}

func (socks5 *Socks5) AuthenticateClient(context *Context) error {
	_, connectionError := context.ClientConnection.Read(context.Chunk[:])
	if connectionError != nil {
		return connectionError
	}
	if context.Chunk[0] != SocksV5 {
		return SocksVersionNotSupported
	}

	clientSupportedMethods := context.Chunk[2 : 2+context.Chunk[1]]

	if socks5.AuthenticationMethod == nil {
		for _, supportedMethod := range clientSupportedMethods {
			if supportedMethod == NoAuthRequired {
				return context.Reply(BasicReply{
					Version:    SocksV5,
					StatusCode: NoAuthRequired,
				})
			}
		}
	} else {
		for _, supportedMethod := range clientSupportedMethods {
			if supportedMethod == UsernamePassword {
				return socks5.UsernamePasswordAuthentication(context)
			}
		}
	}
	_, _ = context.ClientConnection.Write([]byte{SocksV5, NoAcceptableMethods})
	return ClientNotAuthenticated
}
