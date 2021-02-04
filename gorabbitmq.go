package gorabbitmq

import (
	"log"
	"fmt"
	"time"
	"errors"

	"github.com/streadway/amqp"
	"github.com/MrAndreID/gohelpers"
)

type Connection struct {
	Host, Port, Username, Password, VirtualHost string
}

type QueueSetting struct {
	Name string
	Durable, AutoDelete, Exclusive, NoWait bool
	Args map[string]interface{}
}

type ConsumeSetting struct {
	Consumer string
	AutoAck, Exclusive, NoLocal, NoWait bool
	Args map[string]interface{}
}

type OtherSetting struct {
	RoutingKey, Expiration string
	Mandatory, Immediate bool
	Timeout time.Duration
}

func RPCClient(body string, connection Connection, queueSetting QueueSetting, consumeSetting ConsumeSetting, otherSetting OtherSetting) (response string, errorResponse error) {
	url := connection.Host + ":" + connection.Port + "/" + connection.VirtualHost

	log.Println("[AMQP] " + url + " [" + queueSetting.Name + "]")

	dial, err := amqp.Dial("amqp://" + connection.Username + ":" + connection.Password + "@" + url)
	if err != nil {
		gohelpers.ErrorMessage("failed to connect rabbitmq", err)

		response = ""
		errorResponse = err

		return
	} else {
		log.Println("Message => successfully connected rabbitmq.")
	}
	defer dial.Close()

	channel, err := dial.Channel()
	if err != nil {
		gohelpers.ErrorMessage("failed to open a channel in rabbitmq", err)

		response = ""
		errorResponse = err

		return
	} else {
		log.Println("Message => successfully to opened a channel in rabbitmq.")
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		queueSetting.Name,
		queueSetting.Durable,
		queueSetting.AutoDelete,
		queueSetting.Exclusive,
		queueSetting.NoWait,
		queueSetting.Args,
	)
	if err != nil {
		gohelpers.ErrorMessage("failed to declare a queue in rabbitmq", err)

		response = ""
		errorResponse = err

		return
	} else {
		log.Println("Message => successfully to declare a queue in rabbitmq.")
	}

	message, err := channel.Consume(
		queue.Name,
		consumeSetting.Consumer,
		consumeSetting.AutoAck,
		consumeSetting.Exclusive,
		consumeSetting.NoLocal,
		consumeSetting.NoWait,
		consumeSetting.Args,
	)
	if err != nil {
		gohelpers.ErrorMessage("failed to register a consumer in rabbitmq", err)

		response = ""
		errorResponse = err

		return
	} else {
		log.Println("Message => successfully to register a consumer in rabbitmq.")
	}

	corrId := gohelpers.Random("str", 32)

	err = channel.Publish(
		"",
		otherSetting.RoutingKey,
		otherSetting.Mandatory,
		otherSetting.Immediate,
		amqp.Publishing{
			ContentType: "text/plain",
			CorrelationId: corrId,
			ReplyTo: queue.Name,
			Body: []byte(body),
			Expiration: otherSetting.Expiration,
		},
	)
	if err != nil {
		gohelpers.ErrorMessage("failed to publish a message to rabbitmq", err)

		response = ""
		errorResponse = err

		return
	} else {
		log.Println("Message => successfully to publish a message to rabbitmq.")
	}

	flag := make(chan string, 1)

	go func() {
		for data := range message {
			if corrId == data.CorrelationId {
				if body == string(data.Body) {
					gohelpers.ErrorMessage("the rpc server is not responding", nil)
	
					response = ""
					errorResponse = errors.New("the rpc server is not responding.")
				} else {
					response = string(data.Body)
					errorResponse = nil
				}
	
				break
			}
		}

		flag <- "successfully to get a response from rabbitmq."
	}()

	select {
	case result := <-flag:
		log.Println(result)

		fmt.Println()
	
		return
	case <-time.After(otherSetting.Timeout * time.Second):
		gohelpers.ErrorMessage("rpc server responds too long", nil)

		response = ""
		errorResponse = errors.New("rpc server responds too long.")

		fmt.Println()
	
		return
	}
}