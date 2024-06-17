package main

import (
	"context"
	"fmt"
	rabbi "github.com/rabbitmq/amqp091-go"
	"gorik.ko/rabbit/internal"
	"log"
	"time"
)

func main() {
	ctx := context.Background()
	//:::: Rabbit CONN setup

	connection, err := internal.MakeConnection("gorik", "gorik", "localhost:5672", "army")
	if err != nil {
		panic("make rabbit conn :: " + err.Error())
	}

	defer connection.Close()

	//:::: CLIENT setup

	client, err := internal.MakeRabbitClient(connection)
	if err != nil {
		panic("make client ::: " + err.Error())
	}
	defer client.Close()
	//:::: Rabbit CONSUMER conn setup

	consumConnection, err := internal.MakeConnection("gorik", "gorik", "localhost:5672", "army")
	if err != nil {
		panic("make rabbit conn :: " + err.Error())
	}

	defer consumConnection.Close()

	//:::: CLIENT consumer setup

	consumClient, err := internal.MakeRabbitClient(consumConnection)
	if err != nil {
		panic("make client ::: " + err.Error())
	}

	//::: consumer  QUEUE  BINDING setup
	queue, err := consumClient.MakeQueue("", true, false)
	if err != nil {
		panic(err)
	}
	err = consumClient.MakeBinding(queue.Name, queue.Name, "army_callback")
	if err != nil {
		panic(err)
	}
	msgBus, _ := consumClient.Consume(queue.Name, "callbacker", false)
	go func() {

		for msg := range msgBus {

			log.Println("callback msg ::: ", msg.CorrelationId)
		}

	}()

	//:::: SENDING msg
	for i := 0; i < 10; i++ {
		err = client.MakeSend(ctx, "army_callback", "", rabbi.Publishing{
			ContentType:   "text/plain",
			DeliveryMode:  rabbi.Persistent,
			ReplyTo:       queue.Name,
			CorrelationId: fmt.Sprintf("Created msg %d", i),
			Body:          []byte(`fresh...`),
		})
		if err != nil {
			panic(err)
		}

	}

	time.Sleep(time.Second * 15)

}
