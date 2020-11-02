package routes

import (
	"gin/controller"
	"gin/servers/token"
	"gin/ws"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRoutes() *gin.Engine {
	hub := ws.NewHub()
	go hub.Run()
	engine := gin.Default()
	engine.Use(CrossHandler())
	v1 := engine.Group("/api/v1")
	{
		v1.POST("/login", controller.LoginHandler)
		v1.GET("ws", func(context *gin.Context) {
			hub.ServeWs(context)
		})
		authorized := v1.Group("/", token.MiddleTokenAuthHandler)
		{
			authorized.GET("logout", controller.LogoutHandler)
			//v1.GET("ws",primary.Start)

			authorized.GET("home", controller.HomeHandler)
			authorized.GET("room/:room_id", controller.RoomHandler)
			authorized.GET("private-chat", controller.PrivateChatHandler)
			authorized.POST("img-upload", controller.ImageUploadHandler)
			authorized.GET("pagination", controller.PaginationHandler)
		}
	}
	return engine
}

//跨域访问：cross  origin resource share
func CrossHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Origin", "*") // 设置允许访问所有域
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
		context.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma,token,openid,opentoken")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
		context.Header("Access-Control-Max-Age", "172800")
		context.Header("Access-Control-Allow-Credentials", "false")
		context.Set("content-type", "application/json")
		//设置返回格式是json

		if method == "OPTIONS" {
			context.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "OK",
			})
		}

		//处理请求
		context.Next()
	}
}
