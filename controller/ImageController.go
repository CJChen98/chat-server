package controller

import (
	"gin/models"
	"gin/servers/image"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ImageUploadHandler(ctx *gin.Context) {
	kind, ok := ctx.GetQuery("type")

	if !ok {
		ctx.JSON(http.StatusBadRequest, models.JSON{
			Code: 400,
			Msg:  "未设图片类型",
		})
		ctx.Abort()
		return
	}
	switch kind {
	case "user":
		image_server.SaveUserAvatar(ctx, kind)
	case "room":
		image_server.SaveRoomAvatar(ctx, kind)
	case "message":
		image_server.SaveMessageImage(ctx, kind)
	default:
		ctx.JSON(http.StatusBadRequest, models.JSON{
			Code: 400,
			Msg:  "图片类型错误",
		})
		ctx.Abort()
		return
	}
}

func ImageHandler(ctx *gin.Context) {
	kind, ok := ctx.GetQuery("type")

	if !ok {
		ctx.JSON(http.StatusBadRequest, models.JSON{
			Code: 400,
			Msg:  "未设图片类型",
		})
		ctx.Abort()
		return
	}
	switch kind {
	case "user", "room", "message":
		image_server.GetImage(ctx, kind)
	default:
		ctx.JSON(http.StatusBadRequest, models.JSON{
			Code: 400,
			Msg:  "图片类型错误",
		})
		ctx.Abort()
		return
	}
}
