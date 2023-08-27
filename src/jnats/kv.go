package jnats

import "github.com/nats-io/nats.go"

func (u *UserNatsConn) GetKVs() (<-chan string, error) {
	nc, err := u.GetNats()
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	return js.KeyValueStoreNames(), nil
}

func (u *UserNatsConn) AddKV(config *nats.KeyValueConfig) (nats.KeyValue, error) {
	nc, err := u.GetNats()
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	return js.CreateKeyValue(config)
}
