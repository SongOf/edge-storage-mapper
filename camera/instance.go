package camera

import (
	"encoding/json"
	"fmt"
	mappercommon "github.com/SongOf/edge-storage-mapper/common"
	"github.com/SongOf/edge-storage-mapper/globals"
	"github.com/SongOf/edge-storage-mapper/models"
	"github.com/SongOf/edge-storage-mapper/scan"
	"k8s.io/klog/v2"
	"strconv"
	"strings"
	"sync"
	"time"
)

var cameras map[string]*PhysicalCamera
var lock sync.Mutex
var wg sync.WaitGroup

type ReportedCamera struct {
	SerialNumber int
	ValidateCode string
	Ip           string
	Url          string
	State        string
	CreateBy     string
	CreateTime   time.Time
	UpdateBy     string
	UpdateTime   time.Time
	ReportedTime time.Time
}

type PhysicalCamera struct {
	Instance    ReportedCamera
	ReportTimer mappercommon.Timer
	Result      string
	Topic       string
}

func (pc PhysicalCamera) Run() {
	_, ok := scan.HostInfoMap.Load(pc.Instance.Ip)
	if !ok {
		klog.Info("设备掉线，设备ID为：", pc.Instance.SerialNumber)
		pc.Instance.State = mappercommon.OFFLINE
	} else {
		pc.Instance.State = mappercommon.ONLINE
	}
	edgeCamera := models.GetCameraBySerialNumber(pc.Instance.SerialNumber)
	if edgeCamera == nil {
		return
	} else {
		pc.Instance.ValidateCode = edgeCamera.ValidateCode
		pc.Instance.Ip = edgeCamera.Ip
		pc.Instance.Url = edgeCamera.Url
		pc.Instance.CreateBy = edgeCamera.CreateBy
		pc.Instance.CreateTime = edgeCamera.CreateTime
		pc.Instance.UpdateBy = edgeCamera.UpdateBy
		pc.Instance.UpdateTime = edgeCamera.UpdateTime
	}
	pc.Instance.ReportedTime = time.Now()

	pc.Result = fmt.Sprintf("%s", pc.Instance)
	fmt.Println(pc.Result)
	if err := pc.handlerPublish(); err != nil {
		klog.Errorf("publish data to mqtt failed: %v", err)
	}
}

func (pc *PhysicalCamera) handlerPublish() (err error) {
	// construct payload
	var payload []byte
	fmt.Println(pc.Topic)
	if strings.Contains(pc.Topic, "$hw") {
		if payload, err = CreateMessageCameraTwinUpdate(pc.Instance); err != nil {
			klog.Errorf("Create message twin update failed: %v", err)
			return
		}
	} else {
		klog.Info("Not implemented")
		return
	}
	if err = globals.MqttClient.Publish(pc.Topic, payload); err != nil {
		klog.Errorf("Publish topic %v failed, err: %v", pc.Topic, err)
	}

	klog.V(2).Infof("Update value: %s, topic: %s", pc.Result, pc.Topic)
	return
}

func CameraInit() error {
	cameras = make(map[string]*PhysicalCamera)
	bindedDevices := models.GetCameraByMapperId(mappercommon.MAPPER_ID)
	for _, device := range bindedDevices {
		pc := &PhysicalCamera{
			ReportedCamera{
				SerialNumber: device.SerialNumber,
				ValidateCode: device.ValidateCode,
				Url:          device.Url,
				Ip:           device.Ip,
				State:        device.State,
				CreateBy:     device.CreateBy,
				CreateTime:   device.CreateTime,
				UpdateBy:     device.UpdateBy,
				UpdateTime:   device.UpdateTime,
				ReportedTime: time.Now(),
			},
			mappercommon.Timer{
				Function: nil,
				Duration: 60 * 1 * time.Second,
				Times:    0,
				Shutdown: make(chan string),
			},
			"",
			fmt.Sprintf(mappercommon.TopicTwinUpdate, strconv.Itoa(device.SerialNumber)),
		}
		pc.ReportTimer.Function = pc.Run
		cameras[strconv.Itoa(device.SerialNumber)] = pc
	}
	return nil
}

func CameraStart() {
	for _, camera := range cameras {
		go camera.ReportTimer.Start()
	}
}

func CameraAdd(camera ReportedCamera) {
	lock.Lock()
	pc := &PhysicalCamera{
		camera,
		mappercommon.Timer{
			Function: nil,
			Duration: 60 * 1 * time.Second,
			Times:    0,
			Shutdown: make(chan string),
		},
		"",
		fmt.Sprintf(mappercommon.TopicTwinUpdate, strconv.Itoa(camera.SerialNumber)),
	}
	pc.ReportTimer.Function = pc.Run
	cameras[strconv.Itoa(camera.SerialNumber)] = pc
	lock.Unlock()
	go cameras[strconv.Itoa(camera.SerialNumber)].ReportTimer.Start()
}

func CameraDel(serialNumber int) {
	lock.Lock()
	sn := strconv.Itoa(serialNumber)
	camera := cameras[sn]
	camera.ReportTimer.Terminated()
	delete(cameras, sn)
	lock.Unlock()
}

func CreateMessageCameraTwinUpdate(camera ReportedCamera) (msg []byte, err error) {
	var updateMsg mappercommon.DeviceTwinUpdate

	updateMsg.BaseMessage.Timestamp = mappercommon.GetTimestamp()
	updateMsg.Twin = map[string]*mappercommon.MsgTwin{}

	updateMsg.Twin["SerialNumber"] = &mappercommon.MsgTwin{}
	sn := strconv.Itoa(camera.SerialNumber)
	updateMsg.Twin["SerialNumber"].Actual = &mappercommon.TwinValue{Value: &sn}
	updateMsg.Twin["SerialNumber"].Metadata = &mappercommon.TypeMetadata{Type: "int"}

	updateMsg.Twin["ValidateCode"] = &mappercommon.MsgTwin{}
	updateMsg.Twin["ValidateCode"].Actual = &mappercommon.TwinValue{Value: &camera.ValidateCode}
	updateMsg.Twin["ValidateCode"].Metadata = &mappercommon.TypeMetadata{Type: "string"}

	updateMsg.Twin["Ip"] = &mappercommon.MsgTwin{}
	updateMsg.Twin["Ip"].Actual = &mappercommon.TwinValue{Value: &camera.Ip}
	updateMsg.Twin["Ip"].Metadata = &mappercommon.TypeMetadata{Type: "string"}

	updateMsg.Twin["Url"] = &mappercommon.MsgTwin{}
	updateMsg.Twin["Url"].Actual = &mappercommon.TwinValue{Value: &camera.Url}
	updateMsg.Twin["Url"].Metadata = &mappercommon.TypeMetadata{Type: "string"}

	updateMsg.Twin["State"] = &mappercommon.MsgTwin{}
	updateMsg.Twin["State"].Actual = &mappercommon.TwinValue{Value: &camera.State}
	updateMsg.Twin["State"].Metadata = &mappercommon.TypeMetadata{Type: "string"}

	msg, err = json.Marshal(updateMsg)
	return
}
