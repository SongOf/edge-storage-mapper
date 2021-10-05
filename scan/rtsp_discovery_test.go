package scan

import (
	"encoding/json"
	"fmt"
	"k8s.io/klog/v2"
	"testing"
	"time"
)

func TestDiscoveryHosts(t *testing.T) {
	hosts, err := DiscoveryRtspHosts(60 * time.Second)
	if err != nil {
		klog.Error(err)
	}
	for _, host := range hosts {
		data, _ := json.Marshal(host)
		fmt.Println(string(data))
	}
}
