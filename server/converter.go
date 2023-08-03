package server

import (
	"github.com/hiamthach/micro-chat/model"
	"github.com/hiamthach/micro-chat/pb"
)

func convertRoom(room model.Room) *pb.Room {
	return &pb.Room{
		Id:           room.ID,
		RoomSize:     room.RoomSize,
		Participants: room.Participants,
		CreatedBy:    room.CreatedBy,
	}
}

func convertMessage(message model.Message) *pb.Message {
	return &pb.Message{
		Id:        message.ID,
		Content:   message.Content,
		SenderId:  message.SenderID,
		RoomId:    message.RoomID,
		Timestamp: message.Timestamp,
	}
}

func convertMessages(messages []model.Message) []*pb.Message {
	var result []*pb.Message
	for _, message := range messages {
		result = append(result, convertMessage(message))
	}
	return result
}
