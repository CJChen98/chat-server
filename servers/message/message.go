package servers

import (
	"gin/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetMessageList(ctx *gin.Context, room_id string) {
	if len(room_id) != 0 {
		list := models.GetMsgListByRoom(room_id)
		ctx.JSON(http.StatusOK, models.JSON{
			Code: 200,
			Msg:  "共找到 " + strconv.Itoa(len(list)) + "条消息",
			Data: models.Data{
				Messages: list,
			},
		})
	}
}
