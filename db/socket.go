package db

import (
	"fmt"
	"log"

	socketio "github.com/googollee/go-socket.io"
)

var socketServer *socketio.Server

func NewSocketServer() (*socketio.Server, error) {
	if socketServer != nil {
		return socketServer, nil
	} else {
		socketServer = socketio.NewServer(nil)
		socketServer.OnConnect("/", func(s socketio.Conn) error {
			s.SetContext("")
			fmt.Println("connected:", s.ID())
			log.Print("connected:", s.ID())
			return nil
		})

		socketServer.OnEvent("/", "new_message", func(s socketio.Conn, msg string) {
			log.Println("new_message:", msg)
			fmt.Println("new_message:", msg)
			s.Emit("new_message", msg)
		})

	}

	return socketServer, nil
}
