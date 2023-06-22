package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/hiamthach/micro-chat/model"
	pb "github.com/hiamthach/micro-chat/pb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func NewChatServer() *ChatServer {
	return &ChatServer{}
}

type ChatServer struct {
	DB *mongo.Database
	pb.UnimplementedChatServiceServer
}

func (s *ChatServer) Run() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %s", err.Error())
	}

	fmt.Printf("Server listening at %v\n", lis.Addr())

	server := grpc.NewServer()
	pb.RegisterChatServiceServer(server, s)

	return server.Serve(lis)
}

func (s *ChatServer) CreateRoom(ctx context.Context, req *pb.CreateRoomRequest) (*pb.CreateRoomResponse, error) {
	var room model.Room
	room.Name = req.GetName()
	room.Participants = []string{}
	room.CreatedBy = req.GetName()

	// Insert into MongoDB and bind the ID to the room
	result, err := s.DB.Collection("rooms").InsertOne(ctx, &room)
	if err != nil {
		return nil, err
	}

	// Get the inserted ID and assign it to the room object
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("failed to get inserted ID")
	}
	room.ID = insertedID.Hex()

	log.Printf("Created room: %v", room.ID)

	return &pb.CreateRoomResponse{
		Id: room.ID,
	}, nil
}

func (s *ChatServer) JoinRoom(ctx context.Context, req *pb.JoinRoomRequest) (*pb.JoinRoomResponse, error) {
	var room model.Room
	objectID, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		log.Println("Invalid ID")
		return nil, err
	}

	err = s.DB.Collection("rooms").FindOne(ctx, bson.M{"_id": objectID}).Decode(&room)
	if err != nil {
		fmt.Println("Error:", err)
		return &pb.JoinRoomResponse{
			Message: fmt.Sprintf("Room not found: %v", req.GetId()),
		}, nil
	}

	// Create a new struct or map to hold the updated room data
	updateData := bson.M{
		"$push": bson.M{
			"participants": req.GetName(),
		},
	}

	_, err = s.DB.Collection("rooms").UpdateOne(ctx, bson.M{"_id": objectID}, updateData)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return &pb.JoinRoomResponse{
		Message: fmt.Sprintf("Joined room: %v", room.Name),
	}, nil
}
