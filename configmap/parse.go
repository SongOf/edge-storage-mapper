package configmap

import (
	"encoding/json"
	"errors"
	"github.com/SongOf/edge-storage-mapper/common"
	"github.com/SongOf/edge-storage-mapper/globals"
	"io/ioutil"
	"k8s.io/klog/v2"
)

// Parse parse the configmap.
func Parse(path string,
	devices map[string]*globals.HttpDev,
	dms map[string]common.DeviceModel,
	protocols map[string]common.Protocol) error {
	var deviceProfile common.DeviceProfile

	jsonFile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(jsonFile, &deviceProfile); err != nil {
		return err
	}

	for i := 0; i < len(deviceProfile.DeviceInstances); i++ {
		instance := deviceProfile.DeviceInstances[i]
		j := 0
		for j = 0; j < len(deviceProfile.Protocols); j++ {
			if instance.Protocol == deviceProfile.Protocols[j].Name {
				instance.PProtocol = deviceProfile.Protocols[j]
				break
			}
		}
		if j == len(deviceProfile.Protocols) {
			err = errors.New("Protocol not found")
			return err
		}

		if instance.PProtocol.Protocol != "http" {
			continue
		}

		devices[instance.ID] = new(globals.HttpDev)
		devices[instance.ID].Instance = instance
		klog.V(4).Info("Instance: ", instance.ID, instance)
	}

	for i := 0; i < len(deviceProfile.DeviceModels); i++ {
		dms[deviceProfile.DeviceModels[i].Name] = deviceProfile.DeviceModels[i]
	}

	for i := 0; i < len(deviceProfile.Protocols); i++ {
		protocols[deviceProfile.Protocols[i].Name] = deviceProfile.Protocols[i]
	}
	return nil
}
