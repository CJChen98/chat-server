package conversation_server

import (
	"gin/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func FindConversationByUserId(ctx *gin.Context, id string) {
	conversations := models.FindConversationByUserID(id)
	if len(conversations) < 1 {
		ctx.JSON(http.StatusNotFound, models.JSON{
			Code: 404,
			Msg:  "一条历史会话都没有",
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
