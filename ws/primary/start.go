package primary

import (
	"gin/ws"
	"github.com/gin-gonic/gin"
)

func Create() ws.ServeInterface {
	return &ws.GoServe{}
}
func Start(c *gin.Context) {
	Create().RunWs(c)
}
func OnlineUserCount() int {
	return Create().GetOnlineUserCount()
}

func OnlineRoomUserCount(roomId int) int {
	return Create().GetOnlineRoomUserCount(roomId)
}
