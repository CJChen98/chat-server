package servers

import (
	"gin/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetMessageList(ctx *gin.Context, room_id string, page string) {
	if len(room_id) != 0 {
		list, maxpage := models.GetMsgListByRoom(room_id, page)
		if list == nil {
			ctx.JSON(http.StatusOK, models.JSON{
				Code: 404,
				Msg:  strconv.Itoa(maxpage),
			})
			return
		}
		ctx.JSON(http.StatusOK, models.JSON{
			Code: 200,
			Msg:  strconv.Itoa(maxpage),
			Data: models.Data{
				Messages: list,
			},
		})
	}
}
