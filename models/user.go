package models

import (
	"github.com/astaxie/beego/orm"
)

// AddUser to create a new user
func (u *SysUser) AddUser() int64 {
	db := orm.NewOrm()
	pk, err := db.Insert(u)

	if err != nil {
		panic(err)
	}

	return pk
}

// FindUser to find a specific user
func FindUser(username string) SysUser {
	db := orm.NewOrm()
	var user SysUser

	db.QueryTable("sys_user").Filter("login_name", username).One(&user)
	return user
}

// FindUserByID to find a specific user by Id
func FindUserByID(uid int) SysUser {
	db := orm.NewOrm()
	var user SysUser
	db.QueryTable(new(SysUser)).Filter("user_id", uid).One(&user)
	return user
}
