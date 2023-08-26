package jnats

import "github.com/nats-io/nats.go"

func (u *UserNatsConn) GetConsumers(stream string) (<-chan *nats.ConsumerInfo, error) {
	nc, err := u.GetNats()
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	return js.ConsumersInfo(stream), nil
}

func (u *UserNatsConn) AddConsumer(stream string, config *nats.ConsumerConfig) (*nats.ConsumerInfo, error) {
	nc, err := u.GetNats()
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	return js.AddConsumer(stream, config)
}
