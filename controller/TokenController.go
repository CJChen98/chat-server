package controller

import (
	"gin/models"
	"gin/servers/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpdateTokenHandler(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	newToken, err := token.CreateJWT().UpdateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusOK, models.JSON{
			Code: -3,
			Msg:  err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, models.JSON{
			Code: 1,
			Msg:  newToken,
		})
		return
	}
}
