package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

// SysUser struct
type SysUser struct {
	UserId      int    `orm:"auto"`
	UserName    string `form:"user_name"`
	LoginName   string `form:"login_name"`
	Email       string `form:"email"`
	Password    string `form:"password"`
	Phonenumber string `form:"phonenumber"`
	Salt        string `form:"salt"`
	CreateTime  time.Time
	UpdateTime  time.Time
}

// RandomSalt for password
type RandomSalt struct {
	PasswordSalt string
}

// LoginForm for password
type LoginForm struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

// Notes struct
type Notes struct {
	Id        int `orm:"auto"`
	UserId    int
	Comments  string `form:"comments"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type EdgeCamera struct {
	Id             int    `orm:"auto"`
	MapperId       int    `form:"mapper_id"`
	SerialNumber   int    `form:"serial_number"`
	ValidateCode   string `form:"validate_code"`
	Ip             string `form:"ip"`
	Protocol       string `form:"protocol"`
	Url            string `form:"url"`
	State          string `form:"state"`
	Version        int    `form:"version"`
	CreateBy       string `form:"create_by"`
	CreatorId      int    `form:"creator_id"`
	CreateTime     time.Time
	UpdateBy       string `form:"update_by"`
	LastOperatorId int    `form:"last_operator_id"`
	UpdateTime     time.Time
}

func Init() {
	orm.RegisterModel(new(SysUser))
	orm.RegisterModel(new(EdgeCamera))
	orm.RegisterModel(new(Notes))
}
