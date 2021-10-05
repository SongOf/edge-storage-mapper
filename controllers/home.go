package controllers

import (
	"github.com/SongOf/edge-storage-mapper/models"
	"github.com/astaxie/beego"
)

// HomeController define the base controller
type HomeController struct {
	beego.Controller
}

// Home used to get the login view
func (c *HomeController) Home() {
	loggedUserID := c.GetSession("uId")
	loggedUserName := c.GetSession("username")

	if loggedUserID != nil {
		user := models.FindUserByID(loggedUserID.(int))
		c.Layout = "layout.tpl"
		c.LayoutSections = make(map[string]string)
		c.LayoutSections["Nav"] = "navbar.tpl"
		c.LayoutSections["Footer"] = "footer.tpl"
		c.Data["SysUser"] = user
		c.Data["UserName"] = loggedUserName
		c.Data["UserID"] = user.UserId
		c.Data["Title"] = "Welcome"
		c.TplName = "home.tpl"
		c.Render()
	}

	c.Redirect(beego.URLFor("UserController.Login"), 302)
}
