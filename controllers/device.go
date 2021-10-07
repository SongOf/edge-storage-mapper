package controllers

import (
	"fmt"
	mappercommon "github.com/SongOf/edge-storage-mapper/common"
	"github.com/SongOf/edge-storage-mapper/models"
	"github.com/SongOf/edge-storage-mapper/scan"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"k8s.io/klog/v2"
	"strconv"
	"time"
)

type DiscoveryCamera struct {
	Id       int
	Ip       string
	Protocol string
	Port     uint16
	Url      string
	State    string
	BindInfo string
}

type DetailCamera struct {
	SerialNumber int
	ValidateCode string
	Ip           string
	Protocol     string
	Port         uint16
	Url          string
	State        string
	BindInfo     string
	CreateBy     string
	CreateTime   time.Time
	UpdateTime   time.Time
}

type DeviceController struct {
	beego.Controller
}

func (c *DeviceController) DiscoveryList() {
	UserID := c.GetSession("uId")
	UserName := c.GetSession("username")

	if UserID == nil {
		c.Redirect(beego.URLFor("UserController.Login"), 302)
	}
	cameraOnlineMap := scan.HostInfoMap
	bindedCameraMap := map[string]*models.EdgeCamera{}
	cameras := models.GetCameraByMapperId(mappercommon.MAPPER_ID)
	for _, camera := range cameras {
		bindedCameraMap[camera.Ip] = camera
	}
	discoveryCameras := []DiscoveryCamera{}
	cameraOnlineMap.Range(func(key, value interface{}) bool {
		cameraOnline := value.(*scan.HostInfo)
		ports := cameraOnline.Ports
		var tarPort *scan.Port
		for _, port := range ports {
			if port.Name == "rtsp" {
				tarPort = port
			}
		}
		if tarPort == nil {
			return true
		}
		discoveryCamera := DiscoveryCamera{
			Ip:       cameraOnline.Ip,
			Protocol: tarPort.Name,
			Port:     tarPort.Id,
			Url:      "",
			State:    mappercommon.ONLINE,
		}
		bindCamera, ok := bindedCameraMap[cameraOnline.Ip]
		if ok == true {
			discoveryCamera.BindInfo = mappercommon.BINDED
			//"rtsp://admin:YHDYPD@192.168.1.106:554/h264/ch1/main/av_stream"
			discoveryCamera.Url = fmt.Sprintf("rtsp://admin:%s@%s:%s/h264/ch1/main/av_stream", bindCamera.ValidateCode, cameraOnline.Ip, strconv.Itoa(int(discoveryCamera.Port)))
		} else {
			discoveryCamera.BindInfo = mappercommon.UNBIND
		}
		discoveryCameras = append(discoveryCameras, discoveryCamera)
		return true
	})
	c.Data["Cameras"] = discoveryCameras
	c.Data["UserID"] = UserID
	c.Data["UserName"] = UserName
	c.Layout = "layout.tpl"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["Nav"] = "navbar.tpl"
	c.LayoutSections["Footer"] = "footer.tpl"
	c.TplName = "Discovery.tpl"
	c.Data["Title"] = "Edge Device Discovery"
	c.Render()
}

func (c *DeviceController) GetDevice() {
	UserID := c.GetSession("uId")
	UserName := c.GetSession("username")

	ip := c.Ctx.Input.Param(":ip")
	mode := c.Ctx.Input.Param(":mode")
	if UserID == nil {
		c.Redirect(beego.URLFor("UserController.Login"), 302)
	}
	c.Data["UserID"] = UserID
	c.Data["UserName"] = UserName
	c.Layout = "layout.tpl"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["Nav"] = "navbar.tpl"
	c.LayoutSections["Footer"] = "footer.tpl"
	if mode == "bind" {
		cameras := models.GetCameraByMapperId(mappercommon.MAPPER_ID)
		for _, camera := range cameras {
			if camera.Ip == ip {
				errormap := []string{}
				errormap = append(errormap, "该设备已注册\n")
				c.Data["Errors"] = errormap
				c.Layout = "layout.tpl"
				c.LayoutSections = make(map[string]string)
				c.LayoutSections["Nav"] = "navbar.tpl"
				c.LayoutSections["Footer"] = "footer.tpl"
				c.TplName = "bind_device.tpl"
				c.Data["Title"] = "Bind device details"
				c.Data["UserID"] = UserID
				c.Data["UserName"] = UserName
				return
			}
		}
		c.Data["Ip"] = ip
		c.TplName = "bind_device.tpl"
		c.Data["Title"] = "Bind device details"
	} else {
		eCamera := models.GetCameraByIpAndMapperId(mappercommon.MAPPER_ID, ip)
		c.Data["Camera"] = &DetailCamera{
			SerialNumber: eCamera.SerialNumber,
			ValidateCode: eCamera.ValidateCode,
			Ip:           eCamera.Ip,
			Protocol:     eCamera.Protocol,
			Port:         554,
			Url:          eCamera.Url,
			State:        eCamera.State,
			BindInfo:     mappercommon.BINDED,
			CreateBy:     eCamera.CreateBy,
			CreateTime:   eCamera.CreateTime,
			UpdateTime:   eCamera.UpdateTime,
		}
		c.TplName = "view_device.tpl"
		c.Data["Title"] = "View device details"
	}
	c.Render()
}

func (c *DeviceController) BindDevice() {
	UserID := c.GetSession("uId")
	UserName := c.GetSession("username")

	ip := c.Ctx.Input.Param(":ip")
	mode := c.Ctx.Input.Param(":mode")
	if mode != "bind" {
		errormap := []string{}
		errormap = append(errormap, "url错误\n")
		c.Data["Errors"] = errormap
		c.Layout = "layout.tpl"
		c.LayoutSections = make(map[string]string)
		c.LayoutSections["Nav"] = "navbar.tpl"
		c.LayoutSections["Footer"] = "footer.tpl"
		c.TplName = "bind_device.tpl"
		c.Data["Title"] = "Bind device details"
		c.Data["UserID"] = UserID
		c.Data["UserName"] = UserName
		return
	}
	if UserID == nil {
		c.Redirect(beego.URLFor("UserController.Login"), 302)
	}
	serialNumber, err := strconv.Atoi(c.GetString("serialNumber"))
	if err != nil {
		klog.Error("BindDevice error: ", err)
	}
	validateCode := c.GetString("validateCode")
	valid := validation.Validation{}
	valid.Required(serialNumber, "serialNumber")
	valid.Required(validateCode, "validateCode")
	if valid.HasErrors() {
		errormap := []string{}
		for _, err := range valid.Errors {
			errormap = append(errormap, err.Key+" "+err.Message+"\n")
		}
		c.Data["Errors"] = errormap
		c.Layout = "layout.tpl"
		c.LayoutSections = make(map[string]string)
		c.LayoutSections["Nav"] = "navbar.tpl"
		c.LayoutSections["Footer"] = "footer.tpl"
		c.Data["Title"] = "Bind device details"
		c.Data["UserName"] = UserName
		c.TplName = "bind_device.tpl"
		return
	}

	cameraOnlineMap := scan.HostInfoMap
	cameras := models.GetCameraByMapperId(mappercommon.MAPPER_ID)
	for _, camera := range cameras {
		if camera.Ip == ip {
			errormap := []string{}
			errormap = append(errormap, "该设备已注册\n")
			c.Data["Errors"] = errormap
			c.Layout = "layout.tpl"
			c.LayoutSections = make(map[string]string)
			c.LayoutSections["Nav"] = "navbar.tpl"
			c.LayoutSections["Footer"] = "footer.tpl"
			c.TplName = "bind_device.tpl"
			c.Data["Title"] = "Bind device details"
			c.Data["UserID"] = UserID
			c.Data["UserName"] = UserName
			return
		}
	}
	value, ok := cameraOnlineMap.Load(ip)
	if ok != true {
		klog.Error("camera not exists for ip: ", ip)
	}
	tarCamera := value.(*scan.HostInfo)
	ports := tarCamera.Ports
	var tarPort *scan.Port
	for _, port := range ports {
		if port.Name == "rtsp" {
			tarPort = port
		}
	}
	if tarPort == nil {
		klog.Error("tarPort not exists for ip: ", ip)
	}
	edgeDevice := models.EdgeCamera{
		MapperId:       mappercommon.MAPPER_ID,
		SerialNumber:   serialNumber,
		ValidateCode:   validateCode,
		Ip:             tarCamera.Ip,
		Protocol:       tarPort.Name,
		Url:            fmt.Sprintf("rtsp://admin:%s@%s:%s/h264/ch1/main/av_stream", validateCode, tarCamera.Ip, strconv.Itoa(int(tarPort.Id))),
		State:          mappercommon.ONLINE,
		Version:        1,
		CreateBy:       UserName.(string),
		CreatorId:      UserID.(int),
		CreateTime:     time.Now(),
		UpdateBy:       UserName.(string),
		LastOperatorId: UserID.(int),
		UpdateTime:     time.Now(),
	}
	if edgeDevice.AddCamera() != 0 {
		c.Layout = "layout.tpl"
		c.LayoutSections = make(map[string]string)
		c.LayoutSections["Nav"] = "navbar.tpl"
		c.LayoutSections["Footer"] = "footer.tpl"
		c.Data["Title"] = "Bind device details"
		c.Data["UserName"] = UserName
		c.TplName = "bind_device.tpl"
		c.Redirect(beego.URLFor("DeviceController.DiscoveryList"), 302)
	} else {
		c.Data["errorMsg"] = "Something went wrong. Please try again!"
		c.Layout = "layout.tpl"
		c.LayoutSections = make(map[string]string)
		c.LayoutSections["Nav"] = "navbar.tpl"
		c.LayoutSections["Footer"] = "footer.tpl"
		c.Data["Title"] = "Bind device details"
		c.Data["UserName"] = UserName
		c.TplName = "bind_device.tpl"
		c.Render()
	}
}

func (c *DeviceController) DeviceList() {
	UserID := c.GetSession("uId")
	UserName := c.GetSession("username")
	if UserID != nil {
		cameraModels := models.GetCameraByMapperId(mappercommon.MAPPER_ID)
		var cameras []*DetailCamera
		for _, camera := range cameraModels {
			cameras = append(cameras, &DetailCamera{
				SerialNumber: camera.SerialNumber,
				ValidateCode: camera.ValidateCode,
				Ip:           camera.Ip,
				State:        camera.State,
				BindInfo:     mappercommon.BINDED,
			})
		}
		c.Data["Cameras"] = cameras
		c.Layout = "layout.tpl"
		c.LayoutSections = make(map[string]string)
		c.LayoutSections["Nav"] = "navbar.tpl"
		c.LayoutSections["Footer"] = "footer.tpl"
		c.Data["Title"] = "Edge Devices"
		c.Data["UserID"] = UserID
		c.Data["UserName"] = UserName

		c.TplName = "device.tpl"
		c.Render()
	}
	c.Redirect(beego.URLFor("UserController.Login"), 302)
}

func (c *DeviceController) DeleteDevice() {
	UserID := c.GetSession("uId")
	UserName := c.GetSession("username")
	ip := c.Ctx.Input.Param(":ip")
	if UserID == nil {
		c.Redirect(beego.URLFor("UserController.Login"), 302)
	}
	ec := models.EdgeCamera{
		MapperId: mappercommon.MAPPER_ID,
		Ip:       ip,
	}
	if ec.DeleteCamera() != 0 {
		c.Layout = "layout.tpl"
		c.LayoutSections = make(map[string]string)
		c.LayoutSections["Nav"] = "navbar.tpl"
		c.LayoutSections["Footer"] = "footer.tpl"
		c.Data["Title"] = "Edge devices"
		c.Data["UserName"] = UserName
		c.TplName = "device.tpl"
		c.Redirect(beego.URLFor("DeviceController.DeviceList"), 302)
	} else {
		c.Data["errorMsg"] = "Something went wrong. Please try again!"
		c.Layout = "layout.tpl"
		c.LayoutSections = make(map[string]string)
		c.LayoutSections["Nav"] = "navbar.tpl"
		c.LayoutSections["Footer"] = "footer.tpl"
		c.Data["Title"] = "Edge devices"
		c.Data["UserName"] = UserName
		c.TplName = "device.tpl"
		c.Render()
	}
}
