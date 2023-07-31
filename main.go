package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hiamthach/micro-chat/server"
	"github.com/hiamthach/micro-chat/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Can not load config: ", err)
	}

	// Connect to MongoDB
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(config.DBSource).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	dbClient, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = dbClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Send a ping to confirm a successful connection
	if err := dbClient.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to MongoDB!")

	// Connect to Redis
	redisUtil, err := util.NewRedisUtil(config)
	if err != nil {
		log.Fatal("Can not connect to redis: ", err)
	}

	// Connect to gRPC dbClient
	conn, err := grpc.Dial(config.GRPCServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("can not connect to grpc server: %w", err)
	}

	defer conn.Close()

	// run server
	go server.RunGRPCServer(config, dbClient, *redisUtil, conn)
	server.RunGatewayServer(config, dbClient, *redisUtil, conn)
}
