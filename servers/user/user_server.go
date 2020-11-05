package user_service

import (
	"gin/models"
	"gin/servers/security"
	"gin/servers/token"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	}
	if userInfo.ID > 0 {
		//session.SaveAuthSession(c, strconv.Itoa(int(userInfo.ID)))
		tokenString, err := generationToken(c, &userInfo)
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
					User:     u,
					Messages: make([]models.Message, 1),
					Users:    make([]models.User, 1),
				},
			})
			return
		}

	} else {
		c.JSON(http.StatusInternalServerError, models.JSON{Code: http.StatusInternalServerError, Msg: "系统错误"})
		return
	}
}

func Logout(ctx *gin.Context) {

}

func generationToken(ctx *gin.Context, u *models.User) (string, error) {
	claims := token.MyClaims{
		Uid:      u.ID,
		Username: u.Username,
		StandardClaims: jwt.StandardClaims{
			Id:        strconv.Itoa(int(u.ID)),
			NotBefore: time.Now().Unix(),        // 签名生效时间
			ExpiresAt: time.Now().Unix() + 7200, // 过期时间 2小时
			Issuer:    "chitanda-gin-chat",      //签名的发行者
		},
	}
	return token.CreateJWT().CreateToken(claims)
}
