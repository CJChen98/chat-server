package models

import (
	"errors"
	"gin/snow"
)

type Room struct {
	ID           uint   `json:"-"`
	CreatorID    string `json:"creator_id"`
	CreatedAt    uint64 `json:"created_at"`
	RoomName     string `json:"room_name"`
	MemberSize   uint   `json:"member_size"`
	Introduction string `json:"introduction"`
	SnowId       string `json:"id"`
	AvatarPath   string `json:"avatar_path"`
}

func AddRoom(value *Room) (*Room, error) {
	var room Room
	rooms := FindRoomByCreator(value.CreatorID)
	if len(rooms) > 20 {
		return nil, errors.New("每个用户最多只能创建10个群")
	}
	room.CreatorID = value.CreatorID
	room.MemberSize = 1
	room.RoomName = value.RoomName
	room.Introduction = value.Introduction
	room.SnowId = snow.Snowflake.GetStringId()
	err := ChatDB.Create(&room)
	return &room, err.Error
}
func CreateSystemRoom() (*Room, error) {
	var room Room
	ChatDB.Where("creator_id = 1").First(&room)
	if room.ID > 0 {
		return nil, errors.New("系统房间已存在")
	}
	room.CreatorID = " 1"
	room.MemberSize = 1
	room.RoomName = "GinChat"
	room.Introduction = "用户注册默认加入"
	room.SnowId = snow.Snowflake.GetStringId()
	err := ChatDB.Create(&room)
	return &room, err.Error
}
func GetSystemRoomID() string {
	var room Room
	ChatDB.Where("creator_id = 1").First(&room)
	if room.ID < 1 {
		r, _ := CreateSystemRoom()
		return r.SnowId
	}
	return room.SnowId
}
func FindRoomByID(id string) Room {
	var room Room
	ChatDB.Where("snow_id = ? ", id).First(&room)
	return room
}
func FindRoomByCreator(creator string) []Room {
	//ID := strconv.Itoa(int(creator))
	var rooms []Room
	ChatDB.Where("creator_id = ?", creator).Find(&rooms)
	return rooms
}
