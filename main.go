package main

import (
	"encoding/json"
	"github.com/SongOf/edge-storage-mapper/camera"
	mappercommon "github.com/SongOf/edge-storage-mapper/common"
	"github.com/SongOf/edge-storage-mapper/conf"
	"github.com/SongOf/edge-storage-mapper/globals"
	"github.com/SongOf/edge-storage-mapper/models"
	"github.com/SongOf/edge-storage-mapper/routers"
	"github.com/SongOf/edge-storage-mapper/scan"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
	"os"
	"strconv"
	"strings"
	"time"
)

func initViper(confFilePath, configFileFullName string) *viper.Viper {
	oviper := viper.New()
	r := strings.Split(configFileFullName, ".")
	if len(r) != 2 {
		klog.Fatal("invalid config file to init viper")
	}
	oviper.SetConfigName(r[0])
	oviper.SetConfigType(r[1])
	oviper.AddConfigPath(confFilePath)
	err := oviper.ReadInConfig()
	if err != nil {
		klog.Fatal("parse config file error", err)
	}
	return oviper
}

func loadMysqlConfig(v *viper.Viper) *conf.MysqlConf {
	mysqlConf := conf.MysqlConf{}
	if err := v.UnmarshalKey("mysql", &mysqlConf); err != nil {
		klog.Fatal("load database config error", err)
	}
	return &mysqlConf
}

func loadMqttConfig(v *viper.Viper) *conf.MqttConf {
	mqttConf := conf.MqttConf{}
	if err := v.UnmarshalKey("mqtt", &mqttConf); err != nil {
		klog.Fatal("load mqtt config error", err)
	}
	return &mqttConf
}

func init() {
	commonViper := initViper("./conf", "common_config.yaml")

	mysqlConf := loadMysqlConfig(commonViper)
	orm.RegisterDriver("mysql", orm.DRMySQL)
	err := orm.RegisterDataBase("default", "mysql", mysqlConf.User+":"+mysqlConf.Password+"@tcp("+mysqlConf.Host+":"+strconv.Itoa(int(mysqlConf.Port))+")/"+mysqlConf.Database+"?charset=utf8&parseTime=True")
	if err != nil {
		klog.Error("Database connection failed!")
		panic(err)
	} else {
		klog.Info("Database connected successfully!")
	}

	mqttConf := loadMqttConfig(commonViper)
	globals.MqttClient = &mappercommon.MqttClient{
		IP:         mqttConf.Server,
		User:       mqttConf.Username,
		Passwd:     mqttConf.Password,
		Cert:       mqttConf.Certification,
		PrivateKey: mqttConf.PrivateKey,
	}
	if err = globals.MqttClient.Connect(); err != nil {
		klog.Fatal(err)
		os.Exit(1)
	}
}

func main() {

	klog.Info("this mapper ID is: ", mappercommon.MAPPER_ID)
	scanCycle := 60 * 1 * time.Second
	timer := mappercommon.Timer{Function: scan.DiscoveryRtspHosts, Duration: scanCycle, Times: 0}
	go func() {
		timer.Start()
	}()
	go func() {
		for {
			scan.HostInfoMap.Range(func(k, v interface{}) bool {
				s, _ := json.Marshal(v)
				klog.Info(string(s))
				return true
			})
			time.Sleep(60 * time.Second)
		}
	}()
	models.Init()
	routers.Init()

	if err := camera.CameraInit(); err != nil {
		klog.Fatal(err)
		os.Exit(1)
	}
	camera.CameraStart()

	beego.Run()
}
