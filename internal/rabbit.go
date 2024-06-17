package internal

import (
	"context"
	"fmt"
	rabbi "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	Conn    *rabbi.Connection
	Channel *rabbi.Channel
}

func MakeConnection(username, password, host, vhost string) (*rabbi.Connection, error) {
	dial, err := rabbi.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, host, vhost))
	if err != nil {
		return nil, err

	}
	return dial, nil
}

func MakeRabbitClient(conn *rabbi.Connection) (*RabbitClient, error) {
	ch, err := conn.Channel()
	if err != nil {
		return &RabbitClient{}, err
	}
	return &RabbitClient{
		Conn:    conn,
		Channel: ch,
	}, nil
}

func (r *RabbitClient) Close() error {
	return r.Channel.Close()

}

func (r RabbitClient) MakeQueue(name string, durable, autodelete bool) error {
	_, err := r.Channel.QueueDeclare(name, durable, autodelete, false, false, nil)
	return err
}

func (r *RabbitClient) MakeBinding(name, key, exchange string) error {
	return r.Channel.QueueBind(name, key, exchange, false, nil)

}

func (r *RabbitClient) MakeSend(ctx context.Context, exchange, key string, msg rabbi.Publishing) error {
	return r.Channel.PublishWithContext(ctx, exchange, key, true, false, msg)
}
