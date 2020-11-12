package message_server

import (
	"gin/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetMessageList(ctx *gin.Context, conversationId string, page string) {
	if len(conversationId) != 0 {
		list, maxPage := models.GetMsgListByConversationId(conversationId, page)
		if list == nil {
			ctx.JSON(http.StatusOK, models.JSON{
				Code: 201,
				Msg:  strconv.Itoa(maxPage),
			})
			return
		}
		ctx.JSON(http.StatusOK, models.JSON{
			Code: 200,
			Msg:  strconv.Itoa(maxPage),
			Data: models.Data{
				Messages: list,
			},
		})
	}
}
