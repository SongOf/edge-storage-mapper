package conf

import "time"

type MysqlConf struct {
	Host            string
	Port            uint
	User            string
	Password        string
	Database        string
	ConnectTimeout  int
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type MqttConf struct {
	Server        string
	Username      string
	Password      string
	Certification string
	PrivateKey    string
}
