package jnats

import (
	"fmt"
	"github.com/Drakorgaur/jetono-api/src/storage"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	"os"
)

type UserNatsConn struct {
	*storage.AccountServerMap
	*nats.Conn
	creds []byte
}

func (u *UserNatsConn) SetCreds(creds []byte) {
	u.creds = creds
}

func (u *UserNatsConn) UserCredentials() nats.Option {
	userCB := func() (string, error) {
		return nkeys.ParseDecoratedJWT(u.creds)
	}
	sigCB := func(nonce []byte) ([]byte, error) {
		kp, err := nkeys.ParseDecoratedNKey(u.creds)
		if err != nil {
			return nil, fmt.Errorf("unable to extract key pair from file %v", err)
		}
		defer kp.Wipe()

		sig, _ := kp.Sign(nonce)
		return sig, nil
	}
	return nats.UserJWT(userCB, sigCB)
}

// *copy-pasted from nats.go
func sigHandler(nonce []byte, seedFile string) ([]byte, error) {
	contents, err := os.ReadFile(seedFile)
	if err != nil {
		return nil, fmt.Errorf("nats: %v", err)
	}
	kp, err := nkeys.ParseDecoratedNKey(contents)

	defer kp.Wipe()

	sig, _ := kp.Sign(nonce)
	return sig, nil
}

func (u *UserNatsConn) GetNats(options ...nats.Option) (*nats.Conn, error) {
	options = append(options, u.userJWT())
	return nats.Connect(u.ServersList, options...)
}

func (u *UserNatsConn) userJWT() nats.Option {
	return nats.UserJWT(
		// userCB
		func() (string, error) {
			return string(u.creds), nil
		},
		// sigCB
		func(nonce []byte) ([]byte, error) {
			return sigHandler(nonce, string(u.creds))
		},
	)
}
