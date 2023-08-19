package jnats

import (
	"github.com/Drakorgaur/jetono-api/src/storage"
	"github.com/nats-io/nats.go"
	nsc "github.com/nats-io/nsc/cmd"
	"github.com/nats-io/nsc/cmd/store"
)

type UserNatsConn struct {
	*storage.AccountServerMap
	*nats.Conn
	Name  string
	creds []byte
}

func (u *UserNatsConn) SetCreds(creds []byte) {
	u.creds = creds
}

func (u *UserNatsConn) UserCredentials() nats.Option {
	s, err := nsc.GetStoreForOperator(u.Operator)
	if err != nil {
		return nil
	}
	return nats.UserCredentials(s.Resolve(store.Accounts, u.Account, store.Users, u.Name+".jwt"))
}

func (u *UserNatsConn) GetNats(options ...nats.Option) (*nats.Conn, error) {
	options = append(options, u.UserCredentials())
	return nats.Connect(u.ServersList, options...)
}
