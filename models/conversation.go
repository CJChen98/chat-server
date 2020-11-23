package models

import (
	"gorm.io/gorm"
)

type Conversation struct {
	gorm.Model
	ID         uint
	Private    bool   `json:"private"`
	ReceiverID string `json:"receiver_id"`
	UserID     string `json:"user_id"`
}

func AddConversation(private bool, receiver string, user string) (*Conversation, error) {
	var c Conversation
	c.Private = private
	c.ReceiverID = receiver
	c.UserID = user
	err := ChatDB.Create(&c)
	return &c, err.Error
}
func FindConversationByUserID(id string) []Conversation {
	var conversations []Conversation
	ChatDB.Where("user_id = ? OR receiver_id = ?", id, id).Find(&conversations)
	return conversations
}
