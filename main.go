package main

import (
	"encoding/json"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	server, err := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			polling.Default,
			&websocket.Transport{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
		},
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	go server.Serve()

	defer server.Close()

	header := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{""})
	credentials := handlers.AllowCredentials()

	// connect, channel := amqp.ConnectRabbit()

	// amqp.SendRabbitDashboard(server, channel)

	// defer channel.Close()

	// defer connect.Close()

	server.OnConnect("/", func(s socketio.Conn) error {
		log.Print("Socket Connected ID:", s.ID())
		// log.Print("Socket Connected ID:", s.ID(), "\n", s)
		s.Join("socketBC")
		// s.Emit("msg", interface{}(map[string]interface{}{
		// 	"message": "MeowSage",
		// }))
		log.Print("client : ", server.Count())
		return nil
	})

	server.OnError("/", func(s socketio.Conn, err error) {
		log.Print("Socket Error : \n", err)
		log.Print("client : ", server.Count())
		s.Close()
	})

	server.OnDisconnect("/", func(s socketio.Conn, dced string) {
		log.Print("Socket Closed : \n", dced)
		log.Print("\nclient : ", server.Count())
		s.Close()
	})

	r.HandleFunc("/msg", func(w http.ResponseWriter, r *http.Request) {
		var msg interface{}

		err := json.NewDecoder(r.Body).Decode(&msg)
		if err != nil {
			log.Println(err.Error())
		}

		log.Print(msg)

		server.BroadcastToRoom("", "socketBC", "msg", msg)

		json.NewEncoder(w).Encode(msg)

	}).Methods("POST")

	r.Handle("/socket.io/", server)

	log.Print("running on :10801")

	http.ListenAndServe(":10801", handlers.CORS(header, methods, origins, credentials)(r))
}
