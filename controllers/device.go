package controllers

import (
	"fmt"
	mappercommon "github.com/SongOf/edge-storage-mapper/common"
	"github.com/SongOf/edge-storage-mapper/models"
	"github.com/SongOf/edge-storage-mapper/scan"
	"github.com/astaxie/beego"
	"strconv"
)

type DiscoveryCamera struct {
	Ip       string
	Protocol string
	Port     uint16
	Url      string
	State    string
	BindInfo string
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
		cameraOnline := value.(scan.HostInfo)
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

func (c *DeviceController) List() {

}

func (c *DeviceController) Add() {

}

func (c *DeviceController) Delete() {

}
