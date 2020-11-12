package models

import (
	"errors"
)

type Room struct {
	ID           uint   `json:"id"`
	CreatorID    uint   `json:"creator_id"`
	CreatedAt    uint64 `json:"created_at"`
	RoomName     string `json:"room_name"`
	MemberSize   uint   `json:"member_size"`
	Introduction string `json:"introduction"`
	ID_          string `json:"id_"`
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
	err := ChatDB.Create(&room)
	return &room, err.Error
}
func FindRoomByID(id string) Room {
	var room Room
	ChatDB.Where("id = ? ", id).First(&room)
	return room
}
func FindRoomByCreator(creator uint) []Room {
	//id := strconv.Itoa(int(creator))
	var rooms []Room
	ChatDB.Where("creator_id = ?", creator).Find(&rooms)
	return rooms
}
