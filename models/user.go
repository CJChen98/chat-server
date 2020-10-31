package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID         uint
	Username   string    `json:"username" binding:"required,max=16,min=2"`
	Password   string    `json:"password" binding:"required,max=32,min=6"`
	AvatarPath string    `json:"avatar_path"`
	//CreateAt   time.Time `time_format:"2006-01-02 15:04:05"`
	//UpdateAt   time.Time `time_format:"2006-01-02 15:04:05"`
}

func AddUser(value interface{}) (*User, error) {
	var u User
	u.Username = value.(map[string]string)["username"]
	u.Password = value.(map[string]string)["password"]
	//u.AvatarPath = value.(map[string]string)["avatar_path"]
	err:=ChatDB.Create(&u)
	return &u,err.Error
}

func FindUserByField(field, value string) User {
	var u User
	if field == "id" || field == "username" {
		ChatDB.Where(field+" = ? ", value).First(&u)
	}
	return u
}

func SaveAvatarPath(avatarPath string, u User) {
	u.AvatarPath = avatarPath
	ChatDB.Save(&u)
}

func GetOnlineUserList(uids []float64) []map[string]interface{} {
	var list []map[string]interface{}
	ChatDB.Where("id in ?", uids).Find(&list)
	return list
}
