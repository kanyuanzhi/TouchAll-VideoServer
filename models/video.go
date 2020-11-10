package models

type Video struct {
	CameraID int    `json:"camera_id" bson:"camera"`
	Image    []byte `json:"image" bson:"image"`
}
