package model

type Room struct {
	ID           string   `bson:"_id,omitempty"`
	Name         string   `bson:"name"`
	Participants []string `bson:"participants"`
	CreatedBy    string   `bson:"createdBy"`
}
