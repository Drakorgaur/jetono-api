package jnats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"os"
	"path/filepath"
)

type UserNatsConn struct {
	Operator string
	Account  string
	Server   string
	*nats.Conn
	User  string
	creds []byte
}

func (u *UserNatsConn) UserCredentials() nats.Option {
	fmt.Println(u.Operator)
	resolve := filepath.Join(os.Getenv("NKEYS_PATH"), "creds", u.Operator, u.Account, u.User+".creds")
	fmt.Printf("resolve %s\n", resolve)
	return nats.UserCredentials(resolve)
}

func (u *UserNatsConn) GetNats(options ...nats.Option) (*nats.Conn, error) {
	options = append(options, u.UserCredentials())
	return nats.Connect(u.Server, options...)
}
