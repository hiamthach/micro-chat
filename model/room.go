package model

type Room struct {
	ID           string   `bson:"_id,omitempty"`
	RoomSize     uint32   `bson:"room_size"`
	Participants []string `bson:"participants"`
	CreatedBy    string   `bson:"createdBy"`
}
