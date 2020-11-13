package models

type Video struct {
	DataType int
	CameraID int    `json:"camera_id" bson:"camera"`
	Image    []byte `json:"image" bson:"image"`
}
