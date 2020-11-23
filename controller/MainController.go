package controller

import (
	"gin/models"
	conversation_server "gin/servers/conversation"
	"gin/servers/message"
	room_server "gin/servers/room"
	"gin/servers/token"
	"gin/servers/user"
	"gin/ws"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoginHandler(ctx *gin.Context) {
	user_service.Login(ctx)
}
func LogoutHandler(ctx *gin.Context) {
	user_service.Logout(ctx)
}

func FetchHandler(ctx *gin.Context) {
	kind, ok := ctx.GetQuery("type")
	if !ok {
		ctx.JSON(http.StatusBadRequest, models.JSON{
			Code: 444,
			Msg:  "未定义查询类型",
		})
		return
	}
	switch kind {
	case "home":
		id, ok := ctx.GetQuery("id")
		if !ok {
			ctx.JSON(http.StatusBadRequest, models.JSON{
				Code: 444,
				Msg:  "未定义查询ID",
			})
			return
		}
		conversation_server.FindConversationByUserId(ctx, id)
		return
	case "room":
		id, ok := ctx.GetQuery("id")
		if !ok {
			ctx.JSON(http.StatusOK, models.JSON{
				Code: 444,
				Msg:  "未定义查询ID",
			})
			return
		}
		room_server.FindRoomByID(ctx, id)
		return
	case "user":
		id, ok := ctx.GetQuery("id")
		if !ok {
			ctx.JSON(http.StatusOK, models.JSON{
				Code: 444,
				Msg:  "未定义查询ID",
			})
			return
		}
		user_service.FindUserById(ctx, id)
		return
	case "msg":
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
		return
	}
}

func Http2WS(hub *ws.Hub) func(c *gin.Context) {
	return func(ctx *gin.Context) {
		tokenString, _ := ctx.GetQuery("token")
		if tokenString == "" {
			ctx.JSON(http.StatusOK, models.JSON{
				Code: -1,
				Msg:  "tokenString is null !",
			})
			return
		}
		claims, err := token.CreateJWT().ParseToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusOK, models.JSON{
				Code: -2,
				Msg:  err.Error(),
			})
			return
		}
		ctx.Set("userinfo", claims)
		hub.ServeWs(ctx)
	}
}
