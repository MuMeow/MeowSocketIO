package socket

import (
	"log"
)

// Handler Socket
func Handler() {

	// server, err := socketio.NewServer(nil)

	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// server.OnConnect("/", func(s socketio.Conn) error {
	// 	s.SetContext("")
	// 	log.Print("Socket Connected ID:", s.ID())
	// 	// s.Join("bc")
	// 	return nil
	// })

	// return server

	log.Print("Socket Test")
}
