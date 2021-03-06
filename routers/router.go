package routers

import (
	"github.com/SongOf/edge-storage-mapper/controllers"
	"github.com/astaxie/beego"
)

func Init() {
	// SignUp, SignIn & Register related routes
	beego.Router("/", &controllers.UserController{}, "get,post:Login")
	beego.Router("/register", &controllers.UserController{}, "get,post:Register")
	beego.Router("/logout/", &controllers.UserController{}, "get:Logout")

	// Member related routes
	beego.Router("/profile/", &controllers.HomeController{}, "get:Home")

	//Edge cameras related routes
	beego.Router("/discovery/", &controllers.DeviceController{}, "get:DiscoveryList")
	beego.Router("/device/:mode/?:ip", &controllers.DeviceController{}, "get:GetDevice;post:BindDevice")
	beego.Router("/device/", &controllers.DeviceController{}, "get:DeviceList")
	beego.Router("/device/delete/?:ip", &controllers.DeviceController{}, "get:DeleteDevice")
	// User notes related routes
	//beego.Router("/notepad/", &controllers.NotePadController{}, "get,post:CreateNote")
	//beego.Router("/notepad/:mode/:id([0-9]+", &controllers.NotePadController{}, "get:GetNotes;post:UpdateNotes")
	//beego.Router("/notepad/delete/:id([0-9]+", &controllers.NotePadController{}, "get:DeleteNotes")
}
