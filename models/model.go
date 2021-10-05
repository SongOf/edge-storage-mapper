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

func init() {
	orm.RegisterModel(new(SysUser))
	orm.RegisterModel(new(Notes))
}
