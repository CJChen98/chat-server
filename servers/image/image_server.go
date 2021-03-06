package image_server

import (
	"fmt"
	"gin/models"
	"gin/servers/token"
	"gin/snow"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
)

const ImageSavePath = "./tmp/img/"

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
	dts := ImageSavePath + kind + "/"
	checkPath(dts)
	userinfo := ctx.MustGet("userinfo").(*token.MyClaims)
	id := snow.Snowflake.GetStringId()
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
	go deleteOldAvatar(findOldAvatarPath(kind, userinfo.SnowId))
	url := generateImgUrl(kind, id)
	models.SaveUserAvatarPath(url, userinfo.SnowId)
	ctx.JSON(http.StatusOK, models.JSON{
		Code: 200,
		Msg:  url,
	})
}
func SaveRoomAvatar(ctx *gin.Context, kind string) {
	header, err := ctx.FormFile("img")
	id, ok := ctx.GetQuery("imgId")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.JSON{
			Code: 400,
			Msg:  err.Error(),
		})
		ctx.Abort()
		return
	}
	if !ok {
		ctx.JSON(http.StatusBadRequest, models.JSON{
			Code: 400,
			Msg:  "未提供id",
		})
		ctx.Abort()
		return
	}
	dts := ImageSavePath + kind + "/"
	checkPath(dts)
	imgId := snow.Snowflake.GetStringId()
	filepath := path.Join(dts, imgId+".png")
	err = save(header, filepath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.JSON{
			Code: http.StatusInternalServerError,
			Msg:  err.Error(),
		})
		ctx.Abort()
		return
	}
	go deleteOldAvatar(findOldAvatarPath(kind, id))
	url := generateImgUrl(kind, imgId)
	models.SaveRoomAvatarPath(url, imgId)
	ctx.JSON(http.StatusOK, models.JSON{
		Code: 200,
		Msg:  url,
	})
}
func SaveMessageImage(ctx *gin.Context, kind string) {
	header, err := ctx.FormFile("img")
	//id, ok := ctx.GetQuery("id")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.JSON{
			Code: 400,
			Msg:  err.Error(),
		})
		ctx.Abort()
		return
	}
	//if !ok {
	//	if !ok {
	//		ctx.JSON(http.StatusBadRequest, models.JSON{
	//			Code: 400,
	//			Msg:  "未提供图片id",
	//		})
	//		ctx.Abort()
	//		return
	//	}
	//}
	dts := ImageSavePath + kind + "/"
	checkPath(dts)
	id := snow.Snowflake.GetStringId()
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
	//models.SaveMessageImage(url, id)
	ctx.JSON(http.StatusOK, models.JSON{
		Code: 200,
		Msg:  url,
	})
}

func generateImgUrl(kind string, id string) string {
	return "/img?type=" + kind + "&id=" + id
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

func GetImage(ctx *gin.Context, kind string) {
	id, ok := ctx.GetQuery("id")
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
	dts := ImageSavePath + kind + "/" + id + ".png"
	img, err := ioutil.ReadFile(dts)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.JSON{
			Code: 400,
			Msg:  "未找到图片",
		})
		ctx.Abort()
		return
	}
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Header("Content-Disposition", "attachment; filename="+id+".png")
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Accept-Length", fmt.Sprintf("%d", len(img)))
	_, _ = ctx.Writer.Write(img)
}
func findOldAvatarPath(kind string, id string) string {
	root := ImageSavePath + kind + "/"
	var url string
	if kind == "user" {
		url = models.FindUserByField("snow_id", id).AvatarPath
	}
	if kind == "room" {
		url = models.FindRoomByID(id).AvatarPath
	}
	split := strings.Split(url, "=")
	if split[len(split)-1] == "default" {
		return ""
	}
	return root + split[len(split)-1] + ".png"
}
func deleteOldAvatar(path string) {
	info, err := os.Stat(path)
	if err != nil {
		log.Println(err.Error())
		return
	}
	if info.IsDir() {
		log.Println(path + " is dir")
		return
	}
	err = os.Remove(path)
	if err != nil {
		log.Println(path + " 删除失败: " + err.Error())
	} else {
		log.Println(path + " 删除成功")
	}
}
