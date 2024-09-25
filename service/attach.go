package service

import (
	"GinChat/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

func Upload(c *gin.Context) {
	w := c.Writer
	r := c.Request
	// 获取文件
	srcFile, header, err := r.FormFile("file")
	if err != nil {
		utils.RespFail(w, err.Error())
	}
	defer srcFile.Close()

	// 上传
	suffix := ".png"
	ofilName := header.Filename
	tem := strings.Split(ofilName, ".")
	if len(tem) > 1 {
		suffix = "." + tem[len(tem)-1]
	}
	fileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)
	url := "./asset/upload/" + fileName
	detFile, err := os.Create(url)
	if err != nil {
		utils.RespFail(w, err.Error())
	}
	_, err = io.Copy(detFile, srcFile)
	if err != nil {
		utils.RespFail(w, err.Error())
	}
	utils.RespOK(w, url, "发送图片成功")
}
