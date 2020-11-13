package conversation_server

import (
	"gin/models"
	room_server "gin/servers/room"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func FindConversationByUserId(ctx *gin.Context, id string) {
	conversations := models.FindConversationByUserID(id)
	if len(conversations) < 1 {
		conversation := room_server.AddToSystemRoom(id)
		ctx.JSON(http.StatusNotFound, models.JSON{
			Code: 200,
			Msg:  "找到了1条历史会话",
			Data: models.Data{
				Conversations: []models.Conversation{*conversation},
			},
		})
		return
	}
	ctx.JSON(http.StatusOK, models.JSON{
		Code: 200,
		Msg:  "找到了" + strconv.Itoa(len(conversations)) + "条历史会话",
		Data: models.Data{
			Conversations: conversations,
		},
	})
}
