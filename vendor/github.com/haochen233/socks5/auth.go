package socks5

import (
	"bytes"
	"fmt"
	"hash"
	"io"
	"sync"
)

// Authenticator provides socks5's authentication sub negotiation.
type Authenticator interface {
	Authenticate(in io.Reader, out io.Writer) error
}

// NoAuth NO_AUTHENTICATION_REQUIRED implementation.
type NoAuth struct {
}

// Authenticate NO_AUTHENTICATION_REQUIRED Authentication for socks5 Server and Client.
func (n NoAuth) Authenticate(in io.Reader, out io.Writer) error {
	return nil
}

// UserPwdAuth provides socks5 Server Username/Password Authenticator.
type UserPwdAuth struct {
	UserPwdStore
}

// Authenticate provide socks5 Server Username/Password authentication.
func (u UserPwdAuth) Authenticate(in io.Reader, out io.Writer) error {
	uname, passwd, err := u.ReadUserPwd(in)
	if err != nil {
		return err
	}

	err = u.Validate(string(uname), string(passwd))
	if err != nil {
		reply := []byte{1, 1}
		_, err1 := out.Write(reply)
		if err1 != nil {
			return err
		}
		return err
	}

	//authentication successful,then send reply to client
	reply := []byte{1, 0}
	_, err = out.Write(reply)
	if err != nil {
		return err
	}

	return nil
}

// ReadUserPwd read Username/Password request from client
// return username and password.
// Username/Password request format is as follows:
//    +----+------+----------+------+----------+
//    |VER | ULEN |  UNAME   | PLEN |  PASSWD  |
//    +----+------+----------+------+----------+
//    | 1  |  1   | 1 to 255 |  1   | 1 to 255 |
//    +----+------+----------+------+----------+
// For standard details, please see (https://www.rfc-editor.org/rfc/rfc1929.html)
func (u UserPwdAuth) ReadUserPwd(in io.Reader) ([]byte, []byte, error) {

	ulen, err := ReadNBytes(in, 2)
	if err != nil {
		return nil, nil, err
	}

	uname, err := ReadNBytes(in, int(ulen[1]))
	if err != nil {
		return nil, nil, err
	}

	plen, err := ReadNBytes(in, 1)
	if err != nil {
		return nil, nil, err
	}

	passwd := make([]byte, plen[0])
	passwd, err = ReadNBytes(in, int(plen[0]))
	if err != nil {
		return nil, nil, err
	}

	return uname, passwd, nil
}

// UserPwdStore provide username and password storage.
type UserPwdStore interface {
	Set(username string, password string) error
	Del(username string) error
	Validate(username string, password string) error
}

// MemoryStore store username&password in memory.
// the password is encrypt with hash method.
type MemoryStore struct {
	Users map[string][]byte
	mu    sync.Mutex
	hash.Hash
	algoSecret string
}

// NewMemeryStore return a new MemoryStore
func NewMemeryStore(algo hash.Hash, secret string) *MemoryStore {
	return &MemoryStore{
		Users:      make(map[string][]byte),
		Hash:       algo,
		algoSecret: secret,
	}
}

// Set the mapping of username and password.
func (m *MemoryStore) Set(username string, password string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	build := bytes.NewBuffer(nil)
	build.WriteString(password + m.algoSecret)
	cryptPasswd := m.Hash.Sum(build.Bytes())
	m.Users[username] = cryptPasswd
	return nil
}

// UserNotExist the error type used in UserPwdStore.Del() method and
// UserPwdStore.Validate method.
type UserNotExist struct {
	username string
}

func (u UserNotExist) Error() string {
	return fmt.Sprintf("user %s don't exist", u.username)
}

// Del delete by username
func (m *MemoryStore) Del(username string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.Users[username]; !ok {
		return UserNotExist{username: username}
	}

	delete(m.Users, username)
	return nil
}

// Validate validate username and password.
func (m *MemoryStore) Validate(username string, password string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.Users[username]; !ok {
		return UserNotExist{username: username}
	}

	build := bytes.NewBuffer(nil)
	build.WriteString(password + m.algoSecret)
	cryptPasswd := m.Hash.Sum(build.Bytes())
	if !bytes.Equal(cryptPasswd, m.Users[username]) {
		return fmt.Errorf("user %s has bad password", username)
	}
	return nil
}
