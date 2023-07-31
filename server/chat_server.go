package server

import (
	"context"
	"errors"

	"github.com/hiamthach/micro-chat/model"
	"github.com/hiamthach/micro-chat/pb"
	"github.com/hiamthach/micro-chat/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ChatServer struct {
	cache  util.RedisUtil
	config util.Config
	store  *mongo.Client
	pb.UnimplementedChatServiceServer
	clientConn *grpc.ClientConn
}

func NewChatServer(config util.Config, cache util.RedisUtil, store *mongo.Client, conn *grpc.ClientConn) (*ChatServer, error) {
	return &ChatServer{
		cache:      cache,
		config:     config,
		store:      store,
		clientConn: conn,
	}, nil
}

func (s *ChatServer) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.Message, error) {
	roomClient := pb.NewRoomServiceClient(s.clientConn)

	room, err := roomClient.GetRoom(ctx, &pb.RoomId{Id: req.RoomId})
	if err != nil {
		return nil, err
	}

	// Check if username in room
	if isExist := util.Contains(req.SenderId, room.Participants); !isExist {
		return nil, errors.New("user is not in room")
	}

	message := model.Message{
		Content:   req.Content,
		SenderID:  req.SenderId,
		RoomID:    req.RoomId,
		Timestamp: timestamppb.Now(),
	}

	if _, err := s.store.Database("chat-app").Collection("messages").InsertOne(ctx, &message); err != nil {
		return nil, err
	}

	return convertMessage(message), nil
}

func (s *ChatServer) GetMessages(ctx context.Context, req *pb.GetMessagesRequest) (*pb.GetMessagesResponse, error) {
	roomClient := pb.NewRoomServiceClient(s.clientConn)

	room, err := roomClient.GetRoom(ctx, &pb.RoomId{Id: req.RoomId})
	if err != nil {
		return nil, err
	}

	// Check if username in room
	if isExist := util.Contains(req.Username, room.Participants); !isExist {
		return nil, errors.New("user is not in room")
	}

	messages := []model.Message{}

	filter := bson.M{"roomId": req.RoomId}
	cursor, err := s.store.Database("chat-app").Collection("messages").Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	return &pb.GetMessagesResponse{
		Messages: convertMessages(messages),
	}, nil
}
