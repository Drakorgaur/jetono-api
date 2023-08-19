package jnats

import "github.com/nats-io/nats.go"

func (u *UserNatsConn) GetStreams() (<-chan *nats.StreamInfo, error) {
	nc, err := u.GetNats()
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	return js.StreamsInfo(), nil
}
