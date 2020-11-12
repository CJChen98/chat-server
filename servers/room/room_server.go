package room_server

import (
	"gin/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CreateRoom(ctx *gin.Context) {
	creator, _ := strconv.Atoi(ctx.PostForm("creator"))
	name := ctx.PostForm("name")
	introduction := ctx.PostForm("intro")
	room, err := models.AddRoom(&models.Room{
		CreatorID:    uint(creator),
		RoomName:     name,
		Introduction: introduction,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, models.JSON{Code: -1,
			Msg: "创建房间失败:" + err.Error()})
		return
	}
	conversation, err := models.AddConversation(false, room.ID, room.CreatorID)
	if err != nil {
		ctx.JSON(http.StatusOK, models.JSON{Code: -1,
			Msg: "创建房间失败:" + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, models.JSON{Code: 200,
		Msg: "创建房间成功", Data: models.Data{
			Conversations: []models.Conversation{*conversation},
			Rooms:         []models.Room{*room},
		}})
}

func FindRoomByID(ctx *gin.Context, id string) {
	room := models.FindRoomByID(id)
	if room.ID < 1 {
		ctx.JSON(http.StatusNotFound, models.JSON{
			Code: 404,
			Msg:  "群号:" + id + "未找到",
		})
		return
	}
	ctx.JSON(http.StatusOK, models.JSON{
		Code: 2001,
		Msg:  "群号:" + id,
		Data: models.Data{Rooms: []models.Room{room}},
	})
}
