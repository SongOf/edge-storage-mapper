package models

import (
	"github.com/astaxie/beego/orm"
)

// AddNote to create a new note
func (ec *EdgeCamera) AddCamera() int64 {
	db := orm.NewOrm()
	pk, err := db.Insert(ec)

	if err != nil {
		panic(err)
	}

	return pk
}

// GetNotesByUserID to get all the notes of a specific user
func GetCameraBySerialNumber(sn int) *EdgeCamera {
	db := orm.NewOrm()
	var edgeCamera EdgeCamera
	db.QueryTable("edge_camera").Filter("serial_number", sn).One(&edgeCamera)

	return &edgeCamera
}

func GetCameraByMapperId(mId int) []*EdgeCamera {
	db := orm.NewOrm()
	var edgeCameras []*EdgeCamera
	db.QueryTable(new(EdgeCamera)).Filter("mapper_id", mId).All(&edgeCameras)

	return edgeCameras
}

func GetCameraByIpAndMapperId(mId int, ip string) *EdgeCamera {
	db := orm.NewOrm()
	var edgeCamera EdgeCamera
	db.QueryTable("edge_camera").Filter("mapper_id", mId).Filter("ip", ip).One(&edgeCamera)

	return &edgeCamera
}

// GetNotesByUserID to get all the notes of a specific user
func GetCameraAll() []*EdgeCamera {
	db := orm.NewOrm()
	var edgeCameras []*EdgeCamera
	db.QueryTable("edge_camera").All(&edgeCameras)

	return edgeCameras
}

// UpdateNote to update a note
func (ec *EdgeCamera) UpdateCamera() int64 {
	db := orm.NewOrm()
	rowsAffected, err := db.QueryTable("edge_camera").Filter("id", ec.Id).Update(orm.Params{
		"serial_number": ec.SerialNumber,
		"validate_code": ec.ValidateCode,
		"state":         ec.State,
	})

	if err != nil {
		panic(err)
	}

	return rowsAffected
}

// DeleteNote to update a note
func (ec *EdgeCamera) DeleteCamera() int64 {
	db := orm.NewOrm()
	rowsAffected, err := db.QueryTable("edge_camera").Filter("mapper_id", ec.MapperId).Filter("ip", ec.Ip).Delete()

	if err != nil {
		panic(err)
	}

	return rowsAffected
}
