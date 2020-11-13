package models

import (
	"gorm.io/gorm"
	"sort"
	"strconv"
)

type Message struct {
	gorm.Model
	ID             uint
	UserId         string `json:"user_id"`
	Username       string `json:"username"`
	ConversationId int    `json:"conversation_id"`
	ReceiverId     string `json:"receiver_id"`
	Content        string `json:"content"`
	ImageUrl       string `json:"image_url"`
}

func SaveContent(message Message) Message {
	var m Message
	m.Username = message.Username
	m.UserId = message.UserId
	m.Content = message.Content
	m.ConversationId = message.ConversationId
	m.ReceiverId = message.ReceiverId

	if "" != message.ImageUrl {
		m.ImageUrl = message.ImageUrl
	}

	ChatDB.Create(&m)
	return m
}

var pageSize = 20

func GetMsgListByReceiverId(receiverId string, page string) ([]Message, int) {
	var result []Message
	pages, _ := strconv.Atoi(page)
	var count int64
	ChatDB.Model(&Message{}).Where("receiver_id = " + receiverId).Count(&count)
	maxPage := count / int64(pageSize)
	if int64(pages) > maxPage || pages < 0 {
		return nil, int(maxPage)
	}
	offset := (pages) * pageSize
	ChatDB.Model(&Message{}).
		Where("receiver_id = " + receiverId).
		Order("ID desc").
		Offset(offset).
		Limit(pageSize).
		Find(&result)
	return result, int(maxPage)
}
func GetLimitMsg(roomId string, offset int) []map[string]interface{} {

	var results []map[string]interface{}
	ChatDB.Model(&Message{}).
		Select("messages.*, users.username ,users.avatar_id").
		Joins("INNER Join users on users.ID = messages.user_id").
		Where("messages.room_id = " + roomId).
		Where("messages.to_user_id = 0").
		Order("messages.ID desc").
		Offset(offset).
		Limit(100).
		Scan(&results)

	if offset == 0 {
		sort.Slice(results, func(i, j int) bool {
			return results[i]["ID"].(uint32) < results[j]["ID"].(uint32)
		})
	}

	return results
}
