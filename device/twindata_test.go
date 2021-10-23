package device

import (
	"encoding/json"
	"fmt"
	mappercommon "github.com/SongOf/edge-storage-mapper/common"
	"github.com/SongOf/edge-storage-mapper/globals"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"k8s.io/klog/v2"
	"os"
	"sync"
	"testing"
)

func TestTwinData_Run(t *testing.T) {
	globals.MqttClient = &mappercommon.MqttClient{IP: "tcp://192.168.1.103:1883",
		User:       "",
		Passwd:     "",
		Cert:       "",
		PrivateKey: ""}
	if err := globals.MqttClient.Connect(); err != nil {
		klog.Fatal(err)
		os.Exit(1)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	topic := fmt.Sprintf(mappercommon.TopicTwinUpdate, "185339746")
	klog.V(1).Info("Subscribe topic: ", topic)
	globals.MqttClient.Subscribe(topic, printMessage)
	wg.Wait()
}

func printMessage(client mqtt.Client, message mqtt.Message) {
	var dt mappercommon.DeviceTwinUpdate
	err := json.Unmarshal(message.Payload(), &dt)
	if err != nil {
		fmt.Println("Error")
	}
	fmt.Printf("%s", message.Payload())
	fmt.Println()
	//Stat即为属性名称
	fmt.Printf("%s", *dt.Twin["Ip"].Actual.Value)
	fmt.Println()
}
