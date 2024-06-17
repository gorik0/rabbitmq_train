package main

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"
	"gorik.ko/rabbit/internal"
	"log"
	"time"
)

func main() {
	//ctx := context.Background()
	//::: RABBIT conn setup

	connection, err := internal.MakeConnection("gorik", "gorik", "localhost:5672", "army")
	if err != nil {
		panic("make rabbit conn :: " + err.Error())
	}

	defer connection.Close()

	//::: CLIENT setup

	client, err := internal.MakeRabbitClient(connection)
	if err != nil {
		panic("make client ::: " + err.Error())
	}
	defer client.Close()
	//::: RABBIT publish conn setup

	connectionP, err := internal.MakeConnection("gorik", "gorik", "localhost:5672", "army")
	if err != nil {
		panic("make rabbit conn :: " + err.Error())
	}

	defer connectionP.Close()

	//::: CLIENT publish setup

	clientP, err := internal.MakeRabbitClient(connectionP)
	if err != nil {
		panic("make client ::: " + err.Error())
	}
	defer clientP.Close()

	//::: QUEUE setup

	q, err := client.MakeQueue("", true, false)
	if err != nil {
		panic("make queue ::: " + err.Error())
	}
	//::: BINDING setup

	err = client.MakeBinding(q.Name, "", "army_callback")
	if err != nil {
		panic("make binding ::: " + err.Error())
	}
	//::: CONSUMING start

	msgBus, err := client.Consume(q.Name, "vonkomat", false)
	if err != nil {
		return
	}

	//::: CREATE CONTEXT errgroup

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()
	gr, ctx := errgroup.WithContext(ctx)
	gr.SetLimit(10)
	go func() {
		for msg := range msgBus {
			gr.Go(func() error {
				log.Println("Msg ::: ", string(msg.MessageId))
				err := msg.Ack(false)
				if err != nil {
					log.Println("Coudln't ack msg!!!")
					return err
				}
				err = clientP.MakeSend(ctx, "army_callback", msg.ReplyTo, amqp091.Publishing{
					DeliveryMode:  amqp091.Persistent,
					ContentType:   "plain/text",
					Body:          []byte("ACK"),
					CorrelationId: msg.CorrelationId})

				if err != nil {
					return err
				}
				log.Println("Msg ack ", string(msg.MessageId))
				return nil
			})

		}
	}()

	doneCh := make(chan struct{})
	<-doneCh

}
