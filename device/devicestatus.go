package device

import (
	"github.com/SongOf/edge-storage-mapper/common"
	"github.com/SongOf/edge-storage-mapper/driver"
	"github.com/SongOf/edge-storage-mapper/globals"
	"k8s.io/klog/v2"
)

// GetStatus is the timer structure for getting device status.
type GetStatus struct {
	Client *driver.HttpClient
	Status string
	topic  string
}

// Run timer function.
func (gs *GetStatus) Run() {
	gs.Status, _ = gs.Client.GetStatus()
	var payload []byte
	var err error
	if payload, err = common.CreateMessageState(gs.Status); err != nil {
		klog.Errorf("Create message state failed: %v", err)
		return
	}
	if err = globals.MqttClient.Publish(gs.topic, payload); err != nil {
		klog.Errorf("Publish failed: %v", err)
		return
	}
}
