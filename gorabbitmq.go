package gorabbitmq

import (
	"time"
	"errors"
	"sync/atomic"

	"github.com/streadway/amqp"
	"github.com/MrAndreID/golog"
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

type QosSetting struct {
	PrefetchCount, PrefetchSize int
	Global bool
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

type RouteFunc func(string) string

type AMQPConnection struct {
	*amqp.Connection
}

type AMQPChannel struct {
	*amqp.Channel
	closed int32
}

type ActionFunc func(string)

func Client(body string, connection Connection, queueSetting QueueSetting, otherSetting OtherSetting) (errorResponse error) {
	url := connection.Host + ":" + connection.Port + "/" + connection.VirtualHost

	golog.Info("[AMQP - Client] " + url + " [" + queueSetting.Name + "]")

	dial, err := amqp.Dial("amqp://" + connection.Username + ":" + connection.Password + "@" + url)
	if err != nil {
		gohelpers.ErrorMessage("failed to connect rabbitmq", err)

		errorResponse = err

		return
	} else {
		golog.Success("Successfully connected rabbitmq.")
	}
	defer dial.Close()

	channel, err := dial.Channel()
	if err != nil {
		gohelpers.ErrorMessage("failed to open a channel in rabbitmq", err)

		errorResponse = err

		return
	} else {
		golog.Success("Successfully to opened a channel in rabbitmq.")
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

		errorResponse = err

		return
	} else {
		golog.Success("Successfully to declare a queue in rabbitmq.")
	}

	err = channel.Publish(
		"",
		queue.Name,
		otherSetting.Mandatory,
		otherSetting.Immediate,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte(body),
			Expiration: otherSetting.Expiration,
		},
	)
	if err != nil {
		gohelpers.ErrorMessage("failed to publish a message to rabbitmq", err)

		errorResponse = err

		return
	} else {
		golog.Success("Successfully to publish a message to rabbitmq.")

		errorResponse = nil

		return
	}
}

func Server(connection Connection, queueSetting QueueSetting, consumeSetting ConsumeSetting, actionFunc ActionFunc) {
	url := connection.Host + ":" + connection.Port + "/" + connection.VirtualHost

	golog.Info("[AMQP - Server] " + url + " [" + queueSetting.Name + "]")

	dial, err := Dial("amqp://" + connection.Username + ":" + connection.Password + "@" + url)
	if err != nil {
		gohelpers.ErrorMessage("failed to connect rabbitmq [Main]", err)
	} else {
		golog.Success("Successfully connected rabbitmq.")
	}
	defer dial.Close()

	channel, err := dial.Channel()
	if err != nil {
		gohelpers.ErrorMessage("failed to open a channel in rabbitmq [Main]", err)
	} else {
		golog.Success("Successfully to opened a channel in rabbitmq.")
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
		gohelpers.ErrorMessage("failed to declare a queue in rabbitmq [Main]", err)
	} else {
		golog.Success("Successfully to declare a queue in rabbitmq.")
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
		gohelpers.ErrorMessage("failed to register a consumer in rabbitmq [Main]", err)
	} else {
		golog.Success("Successfully to register a consumer in rabbitmq.")
	}

	forever := make(chan bool)

	go func() {
		for data := range message {
			golog.Success("Successfully to consume a message to rabbitmq.")

			actionFunc(string(data.Body))
		}
	}()

	<-forever
}

func RPCClient(body string, connection Connection, queueSetting QueueSetting, consumeSetting ConsumeSetting, otherSetting OtherSetting) (response string, errorResponse error) {
	url := connection.Host + ":" + connection.Port + "/" + connection.VirtualHost

	golog.Info("[AMQP - RPC Client] " + url + " [" + queueSetting.Name + "]")

	dial, err := amqp.Dial("amqp://" + connection.Username + ":" + connection.Password + "@" + url)
	if err != nil {
		gohelpers.ErrorMessage("failed to connect rabbitmq", err)

		response = ""
		errorResponse = err

		return
	} else {
		golog.Success("Successfully connected rabbitmq.")
	}
	defer dial.Close()

	channel, err := dial.Channel()
	if err != nil {
		gohelpers.ErrorMessage("failed to open a channel in rabbitmq", err)

		response = ""
		errorResponse = err

		return
	} else {
		golog.Success("Successfully to opened a channel in rabbitmq.")
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
		golog.Success("Successfully to declare a queue in rabbitmq.")
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
		golog.Success("Successfully to register a consumer in rabbitmq.")
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
		golog.Success("Successfully to publish a message to rabbitmq.")
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

		flag <- "Successfully to get a response from rabbitmq."
	}()

	select {
	case result := <-flag:
		golog.Success(result)
	
		return
	case <-time.After(otherSetting.Timeout * time.Second):
		gohelpers.ErrorMessage("rpc server responds too long", nil)

		response = ""
		errorResponse = errors.New("rpc server responds too long.")
	
		return
	}
}

func RPCServer(connection Connection, queueSetting QueueSetting, qosSetting QosSetting, consumeSetting ConsumeSetting, otherSetting OtherSetting, routeFunc RouteFunc) {
	url := connection.Host + ":" + connection.Port + "/" + connection.VirtualHost

	golog.Info("[AMQP - RPC Server] " + url + " [" + queueSetting.Name + "]")

	dial, err := Dial("amqp://" + connection.Username + ":" + connection.Password + "@" + url)
	if err != nil {
		gohelpers.ErrorMessage("failed to connect rabbitmq [Main]", err)
	} else {
		golog.Success("Successfully connected rabbitmq.")
	}
	defer dial.Close()

	channel, err := dial.Channel()
	if err != nil {
		gohelpers.ErrorMessage("failed to open a channel in rabbitmq [Main]", err)
	} else {
		golog.Success("Successfully to opened a channel in rabbitmq.")
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
		gohelpers.ErrorMessage("failed to declare a queue in rabbitmq [Main]", err)
	} else {
		golog.Success("Successfully to declare a queue in rabbitmq.")
	}

	err = channel.Qos(
		qosSetting.PrefetchCount,
		qosSetting.PrefetchSize,
		qosSetting.Global,
	)
	if err != nil {
		gohelpers.ErrorMessage("failed to set qos in rabbitmq [Main]", err)
	} else {
		golog.Success("Successfully to set qos in rabbitmq.")
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
		gohelpers.ErrorMessage("failed to register a consumer in rabbitmq [Main]", err)
	} else {
		golog.Success("Successfully to register a consumer in rabbitmq.")
	}

	forever := make(chan bool)

	go func() {
		for data := range message {
			response := routeFunc(string(data.Body))

			err = channel.Publish(
				"",
				data.ReplyTo,
				otherSetting.Mandatory,
				otherSetting.Immediate,
				amqp.Publishing{
					ContentType: "text/plain",
					CorrelationId: data.CorrelationId,
					Body: []byte(response),
					Expiration: otherSetting.Expiration,
				},
			)
			if err != nil {
				gohelpers.ErrorMessage("failed to publish a message to rabbitmq [Main]", err)
			} else {
				golog.Success("Successfully to publish a message to rabbitmq.")
			}

			data.Ack(false)
		}
	}()

	<-forever
}

func Dial(url string) (*AMQPConnection, error) {
	dial, err := amqp.Dial(url)
	if err != nil {
		gohelpers.ErrorMessage("failed to connect rabbitmq", err)

		return nil, err
	}

	connection := &AMQPConnection{
		Connection: dial,
	}

	go func() {
		for {
			reason, ok := <-connection.Connection.NotifyClose(make(chan *amqp.Error))
			if !ok {
				gohelpers.ErrorMessage("connection closed", reason)

				break
			}

			for {
				time.Sleep(2 * time.Second)

				dial, err := amqp.Dial(url)
				if err != nil {
					gohelpers.ErrorMessage("failed to reconnect rabbitmq", err)
				} else {
					connection.Connection = dial

					golog.Success("Successfully reconnected rabbitmq.")

					break
				}
			}
		}
	}()

	return connection, nil
}

func (connection *AMQPConnection) Channel() (*AMQPChannel, error) {
	connectionChannel, err := connection.Connection.Channel()
	if err != nil {
		gohelpers.ErrorMessage("failed to open a channel in rabbitmq", err)

		return nil, err
	}

	channel := &AMQPChannel{
		Channel: connectionChannel,
	}

	go func() {
		for {
			reason, ok := <-channel.Channel.NotifyClose(make(chan *amqp.Error))
			if !ok || channel.IsClosed() {
				gohelpers.ErrorMessage("channel closed", reason)

				channel.Close()

				break
			}

			for {
				time.Sleep(2 * time.Second)

				connectionChannel, err := connection.Connection.Channel()
				if err != nil {
					gohelpers.ErrorMessage("failed to recreate rabbitmq channel", err)
				} else {
					channel.Channel = connectionChannel

					golog.Success("Successfully recreated rabbitmq channel.")

					break
				}
			}
		}
	}()

	return channel, nil
}

func (channel *AMQPChannel) Close() error {
	if channel.IsClosed() {
		return amqp.ErrClosed
	}

	atomic.StoreInt32(&channel.closed, 1)

	return channel.Channel.Close()
}

func (channel *AMQPChannel) IsClosed() bool {
	return (atomic.LoadInt32(&channel.closed) == 1)
}

func (channel *AMQPChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	deliveries := make(chan amqp.Delivery)

	go func() {
		for {
			data, err := channel.Channel.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
			if err != nil {
				gohelpers.ErrorMessage("failed to register a consumer in rabbitmq", err)

				time.Sleep(2 * time.Second)

				continue
			}

			for message := range data {
				deliveries <- message
			}

			time.Sleep(2 * time.Second)

			if channel.IsClosed() {
				break
			}
		}
	}()

	return deliveries, nil
}
