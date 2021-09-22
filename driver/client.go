package driver

import (
	"fmt"
	"github.com/SongOf/edge-storage-mapper/common"
	"io/ioutil"
	"net/http"
	"net/url"
)

type HttpConfig struct {
	Addr string
}

// BleClient is the client structure.
type HttpClient struct {
	config HttpConfig
}

func NewClient(config HttpConfig) (hc HttpClient, err error) {
	return HttpClient{config: config}, nil
}

func (hc *HttpClient) GetStat() (string, error) {
	params := url.Values{}
	addr := fmt.Sprintf("http://%s%s", hc.config.Addr, "/system/stat")
	Url, _ := url.Parse(addr)
	params.Set("seconds", "1")
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	fmt.Println(urlPath)
	resp, err := http.Get(urlPath)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func (hc *HttpClient) GetStatus() (string, error) {
	addr := fmt.Sprintf("http://%s%s", hc.config.Addr, "/health")
	Url, _ := url.Parse(addr)
	//如果参数中有中文参数,这个方法会进行URLEncode
	urlPath := Url.String()
	fmt.Println(urlPath)
	resp, err := http.Get(urlPath)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if string(body) == "edge storage agent is alive" {
		return common.DEVSTOK, nil
	}
	return common.DEVSTDISCONN, nil
}
