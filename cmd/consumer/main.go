package main

import (
	"gorik.ko/rabbit/internal"
	"log"
	"strconv"
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
	go func() {
		for msg := range msgBus {

			log.Println("Message recieved ::: ", string(msg.Body), strconv.FormatBool(msg.Redelivered))
			if !msg.Redelivered {
				err := msg.Nack(false, true)
				if err != nil {
					log.Println("Coudln't nack msg!!!")
					continue
				}
			}
			err := msg.Ack(false)
			if err != nil {
				log.Println("Coudln't ack msg!!!")
				continue
			}
		}
	}()

	doneCh := make(chan struct{})
	<-doneCh

}
