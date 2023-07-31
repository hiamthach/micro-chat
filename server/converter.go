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
