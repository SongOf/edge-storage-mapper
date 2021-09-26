package device

import (
	"fmt"
	"github.com/SongOf/edge-storage-mapper/common"
	"github.com/SongOf/edge-storage-mapper/driver"
	"github.com/SongOf/edge-storage-mapper/globals"
	"k8s.io/klog/v2"
	"strings"
)

// TwinData is the timer structure for getting twin/data.
type TwinData struct {
	HttpClient *driver.HttpClient
	Name       string
	Type       string
	Result     string
	Topic      string
}

// Run timer function.
func (td *TwinData) Run() {
	b, err := td.HttpClient.GetStat()
	if err != nil {
		klog.Errorf("Failed to read stat: %s\n", err)
	}

	td.Result = fmt.Sprintf("%s", b)
	fmt.Println(td.Result)
	if err = td.handlerPublish(); err != nil {
		klog.Errorf("publish data to mqtt failed: %v", err)
	}
}

func (td *TwinData) handlerPublish() (err error) {
	// construct payload
	var payload []byte
	fmt.Println(td.Topic)
	if strings.Contains(td.Topic, "$hw") {
		if payload, err = common.CreateMessageTwinUpdate(td.Name, td.Type, td.Result); err != nil {
			klog.Errorf("Create message twin update failed: %v", err)
			return
		}
	} else {
		if payload, err = common.CreateMessageData(td.Name, td.Type, td.Result); err != nil {
			klog.Errorf("Create message data failed: %v", err)
			return
		}
	}
	if err = globals.MqttClient.Publish(td.Topic, payload); err != nil {
		klog.Errorf("Publish topic %v failed, err: %v", td.Topic, err)
	}

	klog.V(2).Infof("Update value: %s, topic: %s", td.Result, td.Topic)
	return
}
