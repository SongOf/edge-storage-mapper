package driver

import (
	"fmt"
	"testing"
)

func TestHttpClient_GetStat(t *testing.T) {
	client, _ := NewClient(HttpConfig{Addr: "127.0.0.1:36000"})
	fmt.Println(client.GetStat())
	fmt.Println(client.GetStatus())
}
