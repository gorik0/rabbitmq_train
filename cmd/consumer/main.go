package main

import (
	"context"
	"golang.org/x/sync/errgroup"
	"gorik.ko/rabbit/internal"
	"log"
	"time"
)

func main() {
	//ctx := context.Background()
	connection, err := internal.MakeConnection("gorik", "gorik", "localhost:5672", "army")
	if err != nil {
		panic("make rabbit conn :: " + err.Error())
	}

	defer connection.Close()

	client, err := internal.MakeRabbitClient(connection)
	if err != nil {
		panic("make client ::: " + err.Error())
	}
	defer client.Close()

	msgBus, err := client.Consume("fresh_blood", "vonkomat", false)
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
				log.Println("Msg ack ", string(msg.MessageId))
				return nil
			})

		}
	}()

	doneCh := make(chan struct{})
	<-doneCh

}
