package server

import (
	"context"
	"errors"

	"github.com/hiamthach/micro-chat/model"
	"github.com/hiamthach/micro-chat/pb"
	"github.com/hiamthach/micro-chat/util"
	pubnub "github.com/pubnub/go/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

type RoomServer struct {
	cache      util.RedisUtil
	config     util.Config
	store      *mongo.Client
	clientConn *grpc.ClientConn
	pn         *pubnub.PubNub
	pb.UnimplementedRoomServiceServer
}

func NewRoomServer(config util.Config, cache util.RedisUtil, store *mongo.Client, conn *grpc.ClientConn, pn *pubnub.PubNub) (*RoomServer, error) {
	return &RoomServer{
		cache:      cache,
		config:     config,
		store:      store,
		clientConn: conn,
		pn:         pn,
	}, nil
}

func (s *RoomServer) GetRoom(ctx context.Context, req *pb.RoomId) (*pb.Room, error) {
	// Convert the RoomId to an ObjectId
	roomID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	// Retrieve the room from the database
	var room model.Room
	filter := bson.M{"_id": roomID}
	if err := s.store.Database("chat-app").Collection("rooms").FindOne(ctx, filter).Decode(&room); err != nil {
		return nil, err
	}

	return convertRoom(room), nil
}

func (s *RoomServer) CreateRoom(ctx context.Context, req *pb.CreateRoomRequest) (*pb.CreateRoomResponse, error) {
	room := model.Room{
		RoomSize:     req.RoomSize,
		CreatedBy:    req.Owner,
		Participants: []string{req.Owner},
	}
	result, err := s.store.Database("chat-app").Collection("rooms").InsertOne(ctx, &room)
	if err != nil {
		return nil, err
	}

	room.ID = result.InsertedID.(primitive.ObjectID).Hex()

	res := &pb.CreateRoomResponse{
		Room: convertRoom(room),
	}

	return res, nil
}

func (s *RoomServer) JoinRoom(ctx context.Context, req *pb.JoinRoomRequest) (*pb.JoinRoomResponse, error) {
	room, err := s.GetRoom(ctx, &pb.RoomId{Id: req.Id})
	if err != nil {
		return nil, err
	}

	roomID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	// Check if the room is full
	if len(room.Participants) >= int(room.RoomSize) {
		return nil, errors.New("room is full")
	}

	// Check if the user is already in the room
	for _, participant := range room.Participants {
		if participant == req.Username {
			return nil, errors.New("user is already in the room")
		}
	}

	// Add the user to the room
	room.Participants = append(room.Participants, req.Username)

	filter := bson.M{"_id": roomID}
	update := bson.M{"$set": bson.M{"participants": room.Participants}}
	if _, err := s.store.Database("chat-app").Collection("rooms").UpdateOne(ctx, filter, update); err != nil {
		return nil, err
	}

	res := &pb.JoinRoomResponse{
		RoomId:   req.Id,
		Username: req.Username,
		Message:  "Joined room successfully",
	}

	return res, nil
}
