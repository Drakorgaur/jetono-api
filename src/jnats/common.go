package jnats

import (
	"fmt"
	"github.com/Drakorgaur/jetono-api/src/storage"
	"github.com/nats-io/nats.go"
	nsc "github.com/nats-io/nsc/cmd"
	"github.com/nats-io/nsc/cmd/store"
)

type UserNatsConn struct {
	*storage.AccountServerMap
	*nats.Conn
	User  string
	creds []byte
}

func (u *UserNatsConn) UserCredentials() nats.Option {
	fmt.Println(u.Operator)
	s, err := nsc.GetStoreForOperator(u.Operator)
	if err != nil {
		return nil
	}
	resolve := s.Resolve(store.Accounts, u.Account, store.Users, u.User+".jwt")
	fmt.Printf("resolve %s\n", resolve)
	return nats.UserCredentials(resolve)
}

func (u *UserNatsConn) GetNats(options ...nats.Option) (*nats.Conn, error) {
	options = append(options, u.UserCredentials())
	return nats.Connect(u.ServersList, options...)
}
