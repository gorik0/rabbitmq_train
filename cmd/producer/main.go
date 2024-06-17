package main

import (
	"context"
	rabbi "github.com/rabbitmq/amqp091-go"
	"gorik.ko/rabbit/internal"
	"time"
)

func main() {
	ctx := context.Background()
	connection, err := internal.MakeConnection("gorik", "gorik", "localhost:5672", "army")
	if err != nil {
		panic("make rabbit conn :: " + err.Error())
	}

	defer connection.Close()

	client, err := internal.MakeRabbitClient(connection)
	if err != nil {
		panic("make client ::: " + err.Error())
	}

	err = client.MakeQueue("fresh_blood", true, false)
	if err != nil {
		panic("queueu::::: " + err.Error())
	}
	err = client.MakeQueue("grandpa_blood", false, true)
	if err != nil {
		panic("queueu::::: " + err.Error())
	}

	err = client.MakeBinding("fresh_blood", "fresh.blood.*", "army_events")
	if err != nil {
		panic("binding ::: " + err.Error())
	}
	err = client.MakeBinding("grandpa_blood", "grandpa.blood.*", "army_events")
	if err != nil {
		panic("binding ::: " + err.Error())
	}

	err = client.MakeSend(ctx, "army_events", "fresh.blood.up", rabbi.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: rabbi.Persistent,
		Body:         []byte(`fresh...`),
	})
	err = client.MakeSend(ctx, "army_events", "grandpa.blood.*", rabbi.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: rabbi.Transient,
		Body:         []byte(`grandPA...`),
	})
	if err != nil {
		panic(err)
	}
	defer client.Close()
	time.Sleep(time.Second * 15)

}
