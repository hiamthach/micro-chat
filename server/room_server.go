package server

import (
	"github.com/hiamthach/micro-chat/pb"
	"github.com/hiamthach/micro-chat/util"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomServer struct {
	cache  util.RedisUtil
	config util.Config
	store  *mongo.Client
	pb.UnimplementedRoomServiceServer
}

func NewRoomServer(config util.Config, cache util.RedisUtil, store *mongo.Client) (*RoomServer, error) {
	return &RoomServer{
		cache:  cache,
		config: config,
		store:  store,
	}, nil
}
