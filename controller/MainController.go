package controller

import (
	user_service "gin/servers/user"
	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	user_service.Login(c)
}
func LogoutHandler(ctx *gin.Context) {
	user_service.Logout(ctx)
}

func HomeHandler(ctx *gin.Context) {

}

func RoomHandler(ctx *gin.Context) {

}

func PrivateChatHandler(ctx *gin.Context) {

}

func PaginationHandler(ctx *gin.Context) {

}
