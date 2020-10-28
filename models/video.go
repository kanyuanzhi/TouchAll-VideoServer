package models

type Video struct {
	Camera int    `json:"camera" bson:"camera"`
	Image  []byte `json:"image" bson:"image"`
}
