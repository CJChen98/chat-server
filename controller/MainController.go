package controller

import (
	"gin/models"
	"gin/servers/message"
	"gin/servers/token"
	"gin/servers/user"
	"gin/ws"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoginHandler(c *gin.Context) {
	user_service.Login(c)
}
func LogoutHandler(ctx *gin.Context) {
	user_service.Logout(ctx)
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
	case "home":
		id, ok := ctx.GetQuery("id")
		if !ok {
			ctx.JSON(http.StatusOK, models.JSON{
				Code: 444,
				Msg:  "未定义查询ID",
			})
			return
		}
		models.FindConversationByUserID(id)
	case "room_msg":
		id, ok := ctx.GetQuery("id")
		if !ok {
			ctx.JSON(http.StatusOK, models.JSON{
				Code: 444,
				Msg:  "未定义查询ID",
			})
			return
		}
		page, ok := ctx.GetQuery("page")
		if !ok {
			page = "0"
		}
		message_server.GetMessageList(ctx, id, page)
	}
}

func Http2WS(hub *ws.Hub) func(c *gin.Context) {
	return func(c *gin.Context) {
		tokenString, _ := c.GetQuery("token")
		if tokenString == "" {
			c.JSON(http.StatusOK, models.JSON{
				Code: -1,
				Msg:  "tokenString is null !",
			})
			return
		}
		_, err := token.CreateJWT().ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusOK, models.JSON{
				Code: -2,
				Msg:  err.Error(),
			})
			return
		}
		hub.ServeWs(c)
	}
}
