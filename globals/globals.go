package globals

import (
	"github.com/SongOf/edge-storage-mapper/common"
	"github.com/SongOf/edge-storage-mapper/driver"
)

type HttpDev struct {
	Instance   common.DeviceInstance
	HttpClient *driver.HttpClient
}

var MqttClient *common.MqttClient
