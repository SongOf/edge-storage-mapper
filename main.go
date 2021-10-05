package main

import (
	"github.com/SongOf/edge-storage-mapper/conf"
	"github.com/SongOf/edge-storage-mapper/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
	"strconv"
	"strings"
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

func init() {
	commonViper := initViper("./conf", "common_config.yaml")

	//appConf := loadAppConfig(commonViper)
	mysqlConf := loadMysqlConfig(commonViper)
	orm.RegisterDriver("mysql", orm.DRMySQL)
	err := orm.RegisterDataBase("default", "mysql", mysqlConf.User+":"+mysqlConf.Password+"@tcp("+mysqlConf.Host+":"+strconv.Itoa(int(mysqlConf.Port))+")/"+mysqlConf.Database+"?charset=utf8&parseTime=True")
	if err != nil {
		klog.Error("Database connection failed!")
		panic(err)
	} else {
		klog.Info("Database connected successfully!")
	}
}

func main() {
	routers.Init()
	beego.Run()
}
