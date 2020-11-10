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
