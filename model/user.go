package model

type User struct {
	ID       string `bson:"_id,omitempty"`
	Username string `bson:"username"`
}
