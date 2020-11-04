package routes

import (
	"gin/controller"
	"gin/models"
	"gin/servers/token"
	"gin/ws"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func InitRoutes() *gin.Engine {
	hub := ws.NewHub()
	go hub.Run()
	engine := gin.Default()
	engine.Use(cors.Default())
	v1 := engine.Group("/api/v1")
	{
		v1.POST("/login", controller.LoginHandler)
		v1.GET("ws/", func(c *gin.Context) {
			tokenString, _ := c.GetQuery("token")
			if tokenString == "" {
				c.JSON(http.StatusOK, models.JSON{
					Code: -1,
					Msg:  "tokenString is null !",
				})
				c.Abort()
				return
			}
			log.Print("get tokenString-->", tokenString)
			_, err := token.CreateJWT().ParseToken(tokenString)
			if err != nil {
				c.JSON(http.StatusOK, models.JSON{
					Code: -2,
					Msg:  err.Error(),
				})
				c.Abort()
				return
			}
			//c.Set("claims", claims)
			hub.ServeWs(c)
		})
		authorized := v1.Group("/", token.MiddleTokenAuthHandler)
		{
			authorized.GET("logout", controller.LogoutHandler)
			authorized.GET("find/", controller.FindHandler)
			authorized.GET("home", controller.HomeHandler)
			authorized.GET("room/:room_id", controller.RoomHandler)
			authorized.GET("private-chat", controller.PrivateChatHandler)
			authorized.POST("img-upload", controller.ImageUploadHandler)
			authorized.GET("pagination", controller.PaginationHandler)
		}
	}
	return engine
}
