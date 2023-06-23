package model

import (
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Message struct {
	ID        string                 `bson:"_id,omitempty"`
	RoomID    string                 `bson:"roomId"`
	SenderID  string                 `bson:"senderId"`
	Content   string                 `bson:"content"`
	Timestamp *timestamppb.Timestamp `bson:"timestamp"`
}
