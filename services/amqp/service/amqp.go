package amqp

import (
	"encoding/json"
	"log"

	socketio "github.com/googollee/go-socket.io"
	"github.com/streadway/amqp"
)

// ConnectRabbit func
func ConnectRabbit() (connect *amqp.Connection, channel *amqp.Channel) {

	connect, err := amqp.Dial("amqp://guest:guest@192.168.1.10:5672")

	if err != nil {
		log.Fatal(err)
	}

	channel, err = connect.Channel()

	err = channel.ExchangeDeclare("test", "topic", false, false, false, false, nil)

	if err != nil {
		log.Fatal(err.Error())
	}

	q, err := channel.QueueDeclare("test", false, false, false, false, nil)

	if err != nil {
		log.Fatal(err.Error())
	}

	err = channel.QueueBind(q.Name, "", "test", false, nil)

	if err != nil {
		log.Fatal(err.Error())
	}

	return connect, channel
}

// SendRabbitDashboard AMQP
func SendRabbitDashboard(server *socketio.Server, channel *amqp.Channel) {

	msg, err := channel.Consume("test", "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	var rb interface{}

	go func() {
		for m := range msg {

			json.Unmarshal(m.Body, &rb)

			server.BroadcastToRoom("", "socketBC", "msg", rb)

			log.Print(rb)
		}
	}()
}
