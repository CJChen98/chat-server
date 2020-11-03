package controller

import (
	"gin/models"
	servers "gin/servers/message"
	user_service "gin/servers/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoginHandler(c *gin.Context) {
	user_service.Login(c)
}
func LogoutHandler(ctx *gin.Context) {
	user_service.Logout(ctx)
}
func HomeHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "NMSL",
	})
}
func FindHandler(ctx *gin.Context) {
	kind, ok := ctx.GetQuery("kind")
	if !ok {
		ctx.JSON(http.StatusOK, models.JSON{
			Code: 444,
			Msg:  "未定义查询类型",
		})
		return
	}
	switch kind {
	case "room_msg":
		id, ok := ctx.GetQuery("id")
		if !ok {
			ctx.JSON(http.StatusOK, models.JSON{
				Code: 444,
				Msg:  "未定义查询ID",
			})
			return
		}
		servers.GetMessageList(ctx, id)
	}
}
func RoomHandler(ctx *gin.Context) {

}

func PrivateChatHandler(ctx *gin.Context) {

}

func PaginationHandler(ctx *gin.Context) {

}
