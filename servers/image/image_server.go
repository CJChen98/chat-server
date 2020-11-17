package image_server

import (
	"gin/models"
	"gin/servers/token"
	"gin/snow"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

func SaveUserAvatar(ctx *gin.Context, kind string) {
	_, header, err := ctx.Request.FormFile("img")
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, models.JSON{
			Code: 400,
			Msg:  err.Error(),
		})
		ctx.Abort()
		return
	}
	log.Println(header.Filename)
	dts := "./tmp/img/" + kind + "/"
	checkPath(dts)
	userinfo := ctx.MustGet("userinfo").(*token.MyClaims)
	filepath := path.Join(dts, userinfo.SnowId+".png")
	err = save(header, filepath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.JSON{
			Code: http.StatusInternalServerError,
			Msg:  err.Error(),
		})
		ctx.Abort()
		return
	}
	url := generateImgUrl(kind, userinfo.SnowId)
	models.SaveUserAvatarPath(url, userinfo.SnowId)
	ctx.JSON(http.StatusOK, models.JSON{
		Code: 200,
		Msg:  "上传成功",
	})
}
func SaveRoomAvatar(ctx *gin.Context, kind string) {
	header, err := ctx.FormFile("img")
	id, ok := ctx.GetQuery("id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.JSON{
			Code: 400,
			Msg:  err.Error(),
		})
		ctx.Abort()
		return
	}
	if !ok {
		if !ok {
			ctx.JSON(http.StatusBadRequest, models.JSON{
				Code: 400,
				Msg:  "未提供图片id",
			})
			ctx.Abort()
			return
		}
	}
	dts := "./tmp/img/" + kind + "/"
	checkPath(dts)
	filepath := path.Join(dts, id+".png")
	err = save(header, filepath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.JSON{
			Code: http.StatusInternalServerError,
			Msg:  err.Error(),
		})
		ctx.Abort()
		return
	}
	url := generateImgUrl(kind, id)
	models.SaveRoomAvatarPath(url, id)
	ctx.JSON(http.StatusOK, models.JSON{
		Code: 200,
		Msg:  "上传成功",
	})
}
func SaveMessageImage(ctx *gin.Context, kind string) {
	header, err := ctx.FormFile("img")
	id, ok := ctx.GetQuery("id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.JSON{
			Code: 400,
			Msg:  err.Error(),
		})
		ctx.Abort()
		return
	}
	if !ok {
		if !ok {
			ctx.JSON(http.StatusBadRequest, models.JSON{
				Code: 400,
				Msg:  "未提供图片id",
			})
			ctx.Abort()
			return
		}
	}
	dts := "./tmp/img/" + kind + "/"
	checkPath(dts)
	//userinfo := ctx.MustGet("userinfo").(*token.MyClaims)
	filepath := path.Join(dts, snow.Snowflake.GetStringId()+".png")
	err = save(header, filepath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.JSON{
			Code: http.StatusInternalServerError,
			Msg:  err.Error(),
		})
		ctx.Abort()
		return
	}
	path := generateImgUrl(kind, id)
	//models.SaveMessageImage(path, id)
	ctx.JSON(http.StatusOK, models.JSON{
		Code: 200,
		Msg:  path,
	})
}

var API_HOST = viper.GetString("app.dev-host")

func generateImgUrl(kind string, id string) string {
	return API_HOST + "/img/?type=" + kind + "&id=" + id
}
func checkPath(dts string) {
	_, err := os.Stat(dts)
	if err != nil {
		if !os.IsExist(err) {
			_ = os.MkdirAll(dts, os.ModePerm)
		}
	}
}
func save(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
