package main

import (
	"github.com/SongOf/edge-storage-mapper/common"
	"github.com/SongOf/edge-storage-mapper/config"
	"github.com/SongOf/edge-storage-mapper/device"
	"github.com/SongOf/edge-storage-mapper/globals"
	"k8s.io/klog/v2"
	"os"
)

func main() {
	var err error
	var c config.Config

	klog.InitFlags(nil)
	defer klog.Flush()

	if err = c.Parse(); err != nil {
		klog.Fatal(err)
		os.Exit(1)
	}
	klog.V(4).Info(c.Configmap)

	globals.MqttClient = &common.MqttClient{IP: c.Mqtt.ServerAddress,
		User:       c.Mqtt.Username,
		Passwd:     c.Mqtt.Password,
		Cert:       c.Mqtt.Cert,
		PrivateKey: c.Mqtt.PrivateKey}
	if err = globals.MqttClient.Connect(); err != nil {
		klog.Fatal(err)
		os.Exit(1)
	}

	if err = device.DevInit(c.Configmap); err != nil {
		klog.Fatal(err)
		os.Exit(1)
	}
	device.DevStart()
}
