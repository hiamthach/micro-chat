package main

import (
	"context"
	"log"

	"github.com/hiamthach/micro-chat/server"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://micro_chat:micro_chat@micro-chat.vmr0qct.mongodb.net/?retryWrites=true&w=majority").SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	chatServer := server.NewChatServer()
	chatServer.DB = client.Database("chat-app")
	if err := chatServer.Run(); err != nil {
		log.Fatalf("failed to run server: %s", err.Error())
	}
}
