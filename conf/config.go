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
