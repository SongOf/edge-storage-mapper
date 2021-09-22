package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {
	config := Config{}
	if err := config.Parse(); err != nil {
		t.Log(err)
		t.FailNow()
	}

	assert.Equal(t, "tcp://127.0.0.1:1883", config.Mqtt.ServerAddress)
	assert.Equal(t, "/opt/kubeedge/deviceProfile.json", config.Configmap)
}
