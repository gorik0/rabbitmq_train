package main

import (
	"context"
	rabbi "github.com/rabbitmq/amqp091-go"
	"gorik.ko/rabbit/internal"
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

	//:::: SENDING msg
	for i := 0; i < 10; i++ {
		err = client.MakeSend(ctx, "army_events", "fresh.blood.up", rabbi.Publishing{
			ContentType:  "text/plain",
			DeliveryMode: rabbi.Persistent,
			Body:         []byte(`fresh...`),
		})
		if err != nil {
			panic(err)
		}

	}

	defer client.Close()
	time.Sleep(time.Second * 15)

}
