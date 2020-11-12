package routes

import (
	"gin/controller"
	"gin/models"
	"gin/servers/token"
	"gin/ws"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRoutes() *gin.Engine {
	engine := gin.Default()
	hub := ws.NewHub()
	go hub.Run()
	engine.Use(CrossHandler())
	v1 := engine.Group("/")
	{
		v1.GET("json", func(context *gin.Context) {
			context.JSON(http.StatusOK, models.JSON{
				Code: 1,
				Msg:  "a",
				Data: models.Data{
					User:          models.User{},
					Users:         make([]models.User, 1),
					Messages:      make([]models.Message, 1),
					Rooms:         make([]models.Room, 1),
					Conversations: make([]models.Conversation, 1),
				},
			})
		})
		v1.POST("login", controller.LoginHandler)
		v1.GET("ws/", controller.Http2WS(hub))
	}
	authorized := v1.Group("/", token.MiddleTokenAuthHandler)
	{
		authorized.GET("/", func(context *gin.Context) {
			context.JSON(http.StatusOK, models.JSON{
				Code: 1,
			})
		})
		authorized.GET("logout", controller.LogoutHandler)
		authorized.GET("fetch/", controller.FetchHandler)
		authorized.POST("create/room", controller.CreateRoomHandler)
		authorized.POST("img-upload", controller.ImageUploadHandler)
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
