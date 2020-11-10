package models

type Request struct {
	RequestType      int `json:"request_type"`
	CameraID         int `json:"camera_id" bson:"camera_id"`
	PreviousCameraID int `json:"previous_camera_id"`
}
