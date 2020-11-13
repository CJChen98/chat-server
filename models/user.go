package models

import (
	"gin/snow"
	"time"
)

type User struct {
	ID         uint      `json:"-"`
	Username   string    `json:"username" binding:"required,max=16,min=2"`
	Password   string    `json:"password" binding:"required,max=32,min=6"`
	AvatarPath string    `json:"avatar_path"`
	SnowId     string    `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at"`
}

func AddUser(value interface{}) (*User, error) {
	var u User
	u.Username = value.(map[string]string)["username"]
	u.Password = value.(map[string]string)["password"]
	u.SnowId = snow.Snowflake.GetStringId()
	//u.AvatarPath = value.(map[string]string)["avatar_path"]
	err := ChatDB.Create(&u)
	u.Password = ""
	return &u, err.Error
}

func FindUserByField(field, value string) User {
	var u User
	if field == "snow_id" || field == "username" {
		ChatDB.Where(field+" = ? ", value).First(&u)
	}
	return u
}

func SaveAvatarPath(avatarPath string, u User) {
	u.AvatarPath = avatarPath
	ChatDB.Save(&u)
}

//func GetOnlineUserList(uids []float64) []map[string]interface{} {
//	var list []map[string]interface{}
//	ChatDB.Where("ID in ?", uids).Find(&list)
//	return list
//}
