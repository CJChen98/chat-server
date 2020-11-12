package models

import (
	"gorm.io/gorm"
	"sort"
	"strconv"
)

type Message struct {
	gorm.Model
	ID             uint
	UserId         int    `json:"user_id"`
	Username       string `json:"username"`
	ConversationId int    `json:"conversation_id"`
	Content        string `json:"content"`
	ImageUrl       string `json:"image_url"`
}

func SaveContent(message Message) Message {
	var m Message
	m.Username = message.Username
	m.UserId = message.UserId
	m.Content = message.Content
	m.ConversationId = message.ConversationId

	if "" != message.ImageUrl {
		m.ImageUrl = message.ImageUrl
	}

	ChatDB.Create(&m)
	return m
}

var pageSize = 20

func GetMsgListByConversationId(conversationId string, page string) ([]Message, int) {
	var result []Message
	pages, _ := strconv.Atoi(page)
	var count int64
	ChatDB.Model(&Message{}).Where("conversation_id = " + conversationId).Count(&count)
	maxPage := count / int64(pageSize)
	if int64(pages) > maxPage || pages < 0 {
		return nil, int(maxPage)
	}
	offset := (pages) * pageSize
	ChatDB.Model(&Message{}).
		Select("user_id", "username", "content", "image_url", "created_at", "conversation_id").
		Where("conversation_id = " + conversationId).
		Order("id desc").
		Offset(offset).
		Limit(pageSize).
		Find(&result)
	return result, int(maxPage)
}
func GetLimitMsg(roomId string, offset int) []map[string]interface{} {

	var results []map[string]interface{}
	ChatDB.Model(&Message{}).
		Select("messages.*, users.username ,users.avatar_id").
		Joins("INNER Join users on users.id = messages.user_id").
		Where("messages.room_id = " + roomId).
		Where("messages.to_user_id = 0").
		Order("messages.id desc").
		Offset(offset).
		Limit(100).
		Scan(&results)

	if offset == 0 {
		sort.Slice(results, func(i, j int) bool {
			return results[i]["id"].(uint32) < results[j]["id"].(uint32)
		})
	}

	return results
}
