package main

import (
	"log"

	"github.com/hiamthach/micro-chat/db"
	"github.com/hiamthach/micro-chat/server"
	"github.com/hiamthach/micro-chat/util"
	"google.golang.org/grpc"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Can not load config: ", err)
	}

	// Connect to MongoDB
	db, err := db.Connect(config.DBSource)
	if err != nil {
		log.Fatal("Can not connect to MongoDB: ", err)
	}

	// Connect to Redis
	redisUtil, err := util.NewRedisUtil(config)
	if err != nil {
		log.Fatal("Can not connect to redis: ", err)
	}

	// Connect to gRPC client
	conn, err := grpc.Dial(config.GRPCServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("can not connect to grpc server: %w", err)
	}

	defer conn.Close()

	// run server
	go server.RunGRPCServer(config, db, *redisUtil, conn)
	server.RunGatewayServer(config, db, *redisUtil, conn)
}
