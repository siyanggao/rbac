package controllers

import (
	"strconv"
	"strings"

	"github.com/astaxie/beego"
)
import (
	"rbac/models"
)

type FormController struct {
	baseController
}

func (c *FormController) Post() {
	user := c.GetSession("user").(models.User)
	result := new(models.BaseResponse)
	f, h, err := c.GetFile("file")
	if err != nil {
		beego.Error("getfile err ", err)
		result.Msg = err.Error()
	}
	defer f.Close()
	beego.Informational(h.Filename)
	var dotIndex int
	if strings.LastIndex(h.Filename, ".") == -1 {
		dotIndex = len(h.Filename)
	} else {
		dotIndex = strings.LastIndex(h.Filename, ".")
	}
	subfix := h.Filename[dotIndex:]
	path := "static/tmp/" + strconv.Itoa(user.Id) + subfix
	err = c.SaveToFile("file", path) // 保存位置在 static/upload, 没有文件夹要先创建
	if err != nil {
		beego.Error(err)
		result.Msg = err.Error()
	} else {
		result.Code = 1
		result.Data = path
	}
	c.Data["json"] = result
	c.ServeJSON()
}
