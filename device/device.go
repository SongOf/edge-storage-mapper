package device

import (
	"encoding/json"
	"fmt"
	mappercommon "github.com/SongOf/edge-storage-mapper/common"
	"github.com/SongOf/edge-storage-mapper/configmap"
	"github.com/SongOf/edge-storage-mapper/driver"
	"github.com/SongOf/edge-storage-mapper/globals"
	"regexp"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"k8s.io/klog/v2"
)

var devices map[string]*globals.HttpDev
var models map[string]mappercommon.DeviceModel
var protocols map[string]mappercommon.Protocol
var wg sync.WaitGroup

// getDeviceID extract the device ID from Mqtt topic.
func getDeviceID(topic string) (id string) {
	re := regexp.MustCompile(`hw/events/device/(.+)/twin/update/delta`)
	return re.FindStringSubmatch(topic)[1]
}

// onMessage callback function of Mqtt subscribe message.
func onMessage(client mqtt.Client, message mqtt.Message) {
	klog.V(2).Info("Receive message", message.Topic())
	// Get device ID and get device instance
	id := getDeviceID(message.Topic())
	if id == "" {
		klog.Error("Wrong topic")
		return
	}
	klog.V(2).Info("Device id: ", id)

	var dev *globals.HttpDev
	var ok bool
	if dev, ok = devices[id]; !ok {
		klog.Error("Device not exist")
		return
	}

	// Get twin map key as the propertyName
	var delta mappercommon.DeviceTwinDelta
	if err := json.Unmarshal(message.Payload(), &delta); err != nil {
		klog.Errorf("Unmarshal message failed: %v", err)
		return
	}
	for twinName, twinValue := range delta.Delta {
		i := 0
		for i = 0; i < len(dev.Instance.Twins); i++ {
			if twinName == dev.Instance.Twins[i].PropertyName {
				break
			}
		}
		if i == len(dev.Instance.Twins) {
			klog.Error("Twin not found: ", twinName)
			continue
		}
		// Desired value is not changed.
		if dev.Instance.Twins[i].Desired.Value == twinValue {
			continue
		}
		dev.Instance.Twins[i].Desired.Value = twinValue
	}
}

// initBLE initialize ble client
func initHttp(protocolConfig configmap.HttpProtocolConfig, name string) (client driver.HttpClient, err error) {
	if protocolConfig.HttpAddress != "" {
		config := driver.HttpConfig{
			Addr: protocolConfig.HttpAddress,
		}
		client, err = driver.NewClient(config)
	}
	return
}

// initTwin initialize the timer to get twin value.
func initTwin(dev *globals.HttpDev) {
	for i := 0; i < len(dev.Instance.Twins); i++ {

		twinData := TwinData{
			HttpClient: dev.HttpClient,
			Name:       dev.Instance.Twins[i].PropertyName,
			Type:       dev.Instance.Twins[i].Desired.Metadatas.Type,
			Topic:      fmt.Sprintf(mappercommon.TopicTwinUpdate, dev.Instance.ID)}
		collectCycle := 10 * time.Second
		timer := mappercommon.Timer{Function: twinData.Run, Duration: collectCycle, Times: 0}
		wg.Add(1)
		go func() {
			defer wg.Done()

			timer.Start()
		}()

	}
}

// initSubscribeMqtt subscribe Mqtt topics.
func initSubscribeMqtt(instanceID string) error {
	topic := fmt.Sprintf(mappercommon.TopicTwinUpdateDelta, instanceID)
	klog.V(1).Info("Subscribe topic: ", topic)
	return globals.MqttClient.Subscribe(topic, onMessage)
}

// initGetStatus start timer to get device status and send to eventbus.
func initGetStatus(dev *globals.HttpDev) {
	getStatus := GetStatus{Client: dev.HttpClient,
		topic: fmt.Sprintf(mappercommon.TopicStateUpdate, dev.Instance.ID)}
	timer := mappercommon.Timer{Function: getStatus.Run, Duration: 10 * time.Second, Times: 0}
	wg.Add(1)
	go func() {
		defer wg.Done()
		timer.Start()
	}()
}

// start start the device.
func start(dev *globals.HttpDev) {
	var protocolConfig configmap.HttpProtocolConfig
	if err := json.Unmarshal([]byte(dev.Instance.PProtocol.ProtocolConfigs), &protocolConfig); err != nil {
		klog.Errorf("Unmarshal ProtocolConfig error: %v", err)
		return
	}

	client, err := initHttp(protocolConfig, protocolConfig.HttpAddress)
	if err != nil {
		klog.Errorf("Init error: %v", err)
		return
	}
	dev.HttpClient = &client

	initTwin(dev)

	if err := initSubscribeMqtt(dev.Instance.ID); err != nil {
		klog.Errorf("Init subscribe mqtt error: %v", err)
		return
	}

	initGetStatus(dev)
}

// DevInit initialize the device datas.
func DevInit(configmapPath string) error {
	devices = make(map[string]*globals.HttpDev)
	models = make(map[string]mappercommon.DeviceModel)
	protocols = make(map[string]mappercommon.Protocol)
	return configmap.Parse(configmapPath, devices, models, protocols)
}

// DevStart start all devices.
func DevStart() {
	for id, dev := range devices {
		klog.V(4).Info("Dev: ", id, dev)
		start(dev)
	}
	wg.Wait()
}
