package onvif

import (
	"fmt"
	"testing"
	"time"
)

func TestStartDiscovery(t *testing.T) {
	deviceList, err := StartDiscovery(60 * time.Second)
	if err != nil {
		t.Error(err)
	}
	js := prettyJSON(&deviceList)
	fmt.Println(js)
}
