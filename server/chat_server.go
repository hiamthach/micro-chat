package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	socketio "github.com/googollee/go-socket.io"
	"github.com/hiamthach/micro-chat/model"
	"github.com/hiamthach/micro-chat/pb"
	"github.com/hiamthach/micro-chat/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ChatServer struct {
	cache      util.RedisUtil
	config     util.Config
	store      *mongo.Client
	clientConn *grpc.ClientConn
	socket     socketio.Server
	pb.UnimplementedChatServiceServer
}

func NewChatServer(config util.Config, cache util.RedisUtil, store *mongo.Client, conn *grpc.ClientConn, socket socketio.Server) (*ChatServer, error) {
	return &ChatServer{
		cache:      cache,
		config:     config,
		store:      store,
		clientConn: conn,
		socket:     socket,
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

	// Clear cache
	redisKey := fmt.Sprintf("messages:%s", req.RoomId)
	if err := s.cache.Clear(ctx, redisKey); err != nil {
		log.Print(err)
	}

	// Emit message to room
	s.socket.BroadcastToRoom("/", req.RoomId, "new_message", convertMessage(message))

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

	redisKey := fmt.Sprintf("messages:%s", req.RoomId)

	// Get data from cache
	cachedResponse, err := s.cache.Get(ctx, redisKey)
	if err == nil {
		// If the response is found in the cache, deserialize it and return
		var response *pb.GetMessagesResponse
		err = json.Unmarshal([]byte(cachedResponse), &response)
		if err != nil {
			return nil, err
		}
		return response, err
	}

	messages := []model.Message{}

	filter := bson.M{"roomId": req.RoomId}
	sortOpts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}})
	cursor, err := s.store.Database("chat-app").Collection("messages").Find(ctx, filter, sortOpts)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	// Set data to cache
	response := &pb.GetMessagesResponse{
		Messages: convertMessages(messages),
	}

	// Serialize the response and store it in Redis cache for future use
	serializedResponse, err := json.Marshal(response)
	if err != nil {
		log.Print(err)
	}
	err = s.cache.Set(ctx, redisKey, serializedResponse, time.Hour)
	if err != nil {
		log.Print(err)
	}

	return response, nil
}
