package models

type Request struct {
	Camera int `json:"camera" bson:"camera"`
}

func NewRequest(camera int) *Request {
	return &Request{
		Camera: camera,
	}
}
