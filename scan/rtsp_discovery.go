package scan

import (
	"fmt"
	"github.com/Ullaakut/nmap"
	"k8s.io/klog/v2"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	rwlock sync.RWMutex
	wg     sync.WaitGroup
)

func parseCIDR() (cidr string, err error) {
	// allow non root user to execute by compare with euid
	if os.Geteuid() != 0 {
		log.Fatal("goscan must run as root.")
	}
	iface = "en0"
	// 初始化 网络信息
	setupNetInfo(iface)
	return ipNet.String(), nil
}

func DiscoveryRtspHosts(duration time.Duration) ([]*HostInfo, error) {
	//allow non root user to execute by compare with euid
	if os.Geteuid() != 0 {
		log.Fatal("goscan must run as root.")
	}
	hostIps := []string{}
	hostInfos := []*HostInfo{}
	cidr, _ := parseCIDR()
	scanner, err := nmap.NewScanner(
		nmap.WithTargets(cidr),
		nmap.WithPingScan(),
	)
	if err != nil {
		return nil, err
	}

	result, _, err := scanner.Run()
	if err != nil {
		return nil, err
	}

	for _, host := range result.Hosts {
		// 查询出所有在线 IP
		ip := fmt.Sprintf("%s", host.Addresses[0])
		if ip != "" {
			// 返回给数组
			hostIps = append(hostIps, ip)
		}
	}

	for _, ip := range hostIps {
		// 遍历每个ip 开启多个 goroutine
		go func(ip string) {
			defer wg.Done()
			host, err := HostsInfo(ip)
			if err != nil {
				klog.Error(err)
			}
			rwlock.RLock()
			if host != nil {
				hostInfos = append(hostInfos, host)
			}
			rwlock.RUnlock()
		}(ip)
		wg.Add(1)
	}
	// 等待所有完成
	wg.Wait()
	res := []*HostInfo{}
	for _, hostInfo := range hostInfos {
		if isRtspHost(hostInfo) {
			res = append(res, hostInfo)
		}
	}
	return res, nil
}

func isRtspHost(hostInfo *HostInfo) bool {
	if hostInfo.Ip == "" {
		return false
	}
	ports := hostInfo.Ports
	for _, port := range ports {
		if port.Name == "rtsp" {
			return true
		}
	}
	return false
}

// 扫描具体信息
func HostsInfo(ips string) (*HostInfo, error) {
	scanner, err := nmap.NewScanner(
		nmap.WithTargets(ips),
		// 开启快速查询  -F
		//nmap.WithFastMode(),
		// 标准查询  -O
		nmap.WithOSDetection(),
	)
	if err != nil {
		return nil, err
	}
	result, _, err := scanner.Run()
	if err != nil {
		return nil, err
	}
	// 初始化结构体
	hosts := new(HostInfo)

	for _, host := range result.Hosts {
		// 过滤 主机 条件
		for _, match := range host.OS.Matches {
			os_name := match.Name
			if strings.Contains(os_name, "Linux") && !strings.Contains(os_name, "Android") {
				rwlock.Lock()
				hosts.OsName = match.Name
				hosts.Ip = ips
				rwlock.Unlock()
				//    查主机  端口 和服务 信息
				for _, port := range host.Ports {
					rwlock.Lock()
					hosts.Ports = append(hosts.Ports, &Port{
						Id:    port.ID,
						State: port.State.State,
						Name:  port.Service.Name,
					})
					rwlock.Unlock()
				}
			}
		}
	}
	return hosts, nil
}
