package model

type Message struct {
	ID        string `bson:"_id,omitempty"`
	RoomID    string `bson:"roomId"`
	SenderID  string `bson:"senderId"`
	Content   string `bson:"content"`
	Timestamp int64  `bson:"timestamp"`
}
