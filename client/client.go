package main

import (
	"context"
	"log"
	"time"

	pb "github.com/hiamthach/micro-chat/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := pb.NewChatServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// params := &pb.CreateRoomRequest{
	// 	Name: "Room",
	// }
	// res, err := c.CreateRoom(ctx, params)
	// if err != nil {
	// 	panic(err)
	// }

	// join := &pb.JoinRoomRequest{
	// 	Id:   "64942cbcff10f89222a95c3e",
	// 	Name: "Thach",
	// }

	// msg, err := c.JoinRoom(ctx, join)
	// if err != nil {
	// 	panic(err)
	// }

	// log.Print(msg)

	message := &pb.SendMessageRequest{
		RoomId:    "64942cbcff10f89222a95c3e",
		SenderId:  "Thach",
		Content:   "Hello",
		Timestamp: timestamppb.New(time.Now()),
	}

	mess, err := c.SendMessage(ctx, message)
	if err != nil {
		panic(err)
	}

	log.Print(mess)

	req := &pb.GetRoomMessagesRequest{
		RoomId: "64942cbcff10f89222a95c3e",
	}

	messages, err := c.GetRoomMessages(ctx, req)
	if err != nil {
		panic(err)
	}
	log.Print(messages.Messages)
}
