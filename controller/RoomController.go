package controller

import (
	room_server "gin/servers/room"
	"github.com/gin-gonic/gin"
)

func CreateRoomHandler(c *gin.Context) {
	room_server.CreateRoom(c)
}
