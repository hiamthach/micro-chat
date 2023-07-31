package db

import (
	"fmt"

	socketio "github.com/googollee/go-socket.io"
)

func NewSocketServer() (*socketio.Server, error) {
	socketServer := socketio.NewServer(nil)
	socketServer.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	return socketServer, nil
}
