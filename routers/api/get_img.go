package api

import (
	"github.com/gin-gonic/gin"
)

// GetImg 请求url /api/getimg?imgname=... 返回图片
func GetImg(c *gin.Context) {
	imgname := c.Query("imgname")
	//log.Println(imgname)

	// 返回文件写法1
	//file, _ := ioutil.ReadFile("./pic/"+imgname)
	//_, err := c.Writer.WriteString(string(file))
	//if err != nil {
	//	return
	//}

	// 返回文件写法2
	c.File("./sticker/" + imgname)
}
