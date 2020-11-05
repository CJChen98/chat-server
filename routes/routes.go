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
	//config := cors.DefaultConfig()
	//config.AllowAllOrigins = true
	//engine.Use(cors.New(cors.Config{
	//	//AllowAllOrigins: true,
	//	AllowOriginFunc:  func(origin string) bool { return true },
	//	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
	//	AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
	//	AllowCredentials: true,
	//	MaxAge:           12 * time.Hour,
	//}))

	engine.Use(CrossHandler())
	v1 := engine.Group("/")
	{
		v1.GET("/", func(context *gin.Context) {
			_, _ = context.Writer.WriteString("hello")
		})
		v1.POST("login", controller.LoginHandler)
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