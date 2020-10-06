package amqp

import (
	"encoding/json"
	"log"
	"os"

	socketio "github.com/googollee/go-socket.io"
	"github.com/streadway/amqp"
)

// ConnectRabbit func
func ConnectRabbit() (connect *amqp.Connection, channel *amqp.Channel) {

	connect, err := amqp.Dial(os.Getenv("RABBIT_ADDRESS"))

	if err != nil {
		log.Fatal(err)
	}

	channel, err = connect.Channel()

	err = channel.ExchangeDeclare(os.Getenv("RABBIT_EXCHANGE"), "topic", false, false, false, false, nil)

	if err != nil {
		log.Fatal(err.Error())
	}

	q, err := channel.QueueDeclare(os.Getenv("RABBIT_QUEUE"), false, false, false, false, nil)

	if err != nil {
		log.Fatal(err.Error())
	}

	err = channel.QueueBind(q.Name, "", os.Getenv("RABBIT_EXCHANGE"), false, nil)

	if err != nil {
		log.Fatal(err.Error())
	}

	return connect, channel
}

// SendRabbitDashboard AMQP
func SendRabbitDashboard(server *socketio.Server, channel *amqp.Channel) {

	msg, err := channel.Consume(os.Getenv("RABBIT_QUEUE"), "", true, false, false, false, nil)
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
