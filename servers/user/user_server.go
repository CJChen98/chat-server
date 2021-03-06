package user_service

import (
	"gin/models"
	room_server "gin/servers/room"
	"gin/servers/security"
	"gin/servers/token"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Login(c *gin.Context) {
	username := c.PostForm("username")
	pwd := c.PostForm("password")
	u := models.User{
		Username: username,
		Password: pwd,
	}
	if err := c.ShouldBind(&u); err != nil {
		c.JSON(http.StatusBadRequest, models.JSON{Code: 500, Msg: "用户名或密码格式不规范"})
		return
	}
	userInfo := models.FindUserByField("username", username)
	md5Pwd := security.Md5Encrypt(pwd)
	if userInfo.ID > 0 {
		if userInfo.Password != md5Pwd {
			c.JSON(http.StatusBadRequest, models.JSON{Code: http.StatusBadRequest, Msg: "密码错误"})
			return
		}
	} else {
		newUser, _ := models.AddUser(map[string]string{
			"username": username,
			"password": md5Pwd,
		})
		userInfo = *newUser
		_ = room_server.AddToSystemRoom(userInfo.SnowId)
	}
	if userInfo.ID > 0 {
		tokenString, err := generationToken(&userInfo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.JSON{Code: http.StatusInternalServerError, Msg: "create tokenString failure"})
			return
		} else {
			u := userInfo
			u.Password = ""
			c.JSON(http.StatusOK, models.JSON{
				Code:  200,
				Msg:   "登录成功",
				Token: tokenString,
				Data: models.Data{
					Users: []models.User{u},
				},
			})
			return
		}

	} else {
		c.JSON(http.StatusInternalServerError, models.JSON{Code: http.StatusInternalServerError, Msg: "系统错误"})
		return
	}
}
func FindUserById(ctx *gin.Context, id string) {
	u := models.FindUserByField("snow_id", id)
	if u.ID < 1 {
		ctx.JSON(http.StatusNotFound, models.JSON{
			Code: 404,
			Msg:  "用户:" + id + "未找到",
		})
		return
	}
	ctx.JSON(http.StatusOK, models.JSON{
		Code: 2002,
		Msg:  "用户:" + id,
		Data: models.Data{Users: []models.User{u}},
	})
	return
}
func Logout(ctx *gin.Context) {

}

func generationToken(u *models.User) (string, error) {
	claims := token.MyClaims{
		User: *u,
		StandardClaims: jwt.StandardClaims{
			Id:        u.SnowId,
			NotBefore: time.Now().Unix(),        // 签名生效时间
			ExpiresAt: time.Now().Unix() + 7200, // 过期时间 2小时
			Issuer:    "chitanda-gin-chat",      //签名的发行者
		},
	}
	return token.CreateJWT().CreateToken(claims)
}
