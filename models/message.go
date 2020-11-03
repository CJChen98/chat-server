package models

import (
	"gorm.io/gorm"
	"sort"
)

type Message struct {
	gorm.Model
	ID       uint
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
	ToUserId int    `json:"to_user_id"`
	RoomId   int    `json:"room_id"`
	Content  string `json:"content"`
	ImageUrl string `json:"image_url"`
}

func SaveContent(message Message) Message {
	var m Message
	m.Username = message.Username
	m.UserId = message.UserId
	m.ToUserId = message.ToUserId
	m.Content = message.Content
	m.RoomId = message.RoomId

	if "" != message.ImageUrl {
		m.ImageUrl = message.ImageUrl
	}

	ChatDB.Create(&m)
	return m
}
func GetMsgListByRoom(roomId string) []Message {
	var result []Message
	ChatDB.Model(&Message{}).
		Select("user_id", "username", "content", "image_url", "created_at").
		Where("room_id = " + roomId).Order("id desc").Find(&result)
	return result
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

func GetLimitPrivateMsg(uid, toUId string, offset int) []map[string]interface{} {

	var results []map[string]interface{}
	ChatDB.Model(&Message{}).
		Select("messages.*, users.username ,users.avatar_id").
		Joins("INNER Join users on users.id = messages.user_id").
		Where("(" +
			"(" + "messages.user_id = " + uid + " and messages.to_user_id=" + toUId + ")" +
			" or " +
			"(" + "messages.user_id = " + toUId + " and messages.to_user_id=" + uid + ")" +
			")").
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
