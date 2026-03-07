package dto

type UserEvent struct {
	UserID int32     `json:"UserId"`
	Event  EventType `json:"Event"`
	Error  string    `json:"Error"`
}

type EventType string

const (
	Add    EventType = "ADD"
	Delete EventType = "DELETE"
	Update  EventType = "UPDATE"
	GetAll EventType = "GET_ALL"
)
