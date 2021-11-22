package models

import (
	"bryson.foundation/kbuildresource/common"
	"github.com/astaxie/beego/orm"
)

type Request struct {
	ID int `json:"id"orm:"column(id)"`
	Name string `json:"name"orm:"column(name)"`
	Message string `json:"message"orm:"column(message)"`
	Status string `json:"status"orm:"column(status)"`
	RequestType string `json:"requestType"orm:"column(request_type)"`
	RequestDTO string `json:"request_dto"orm:"column(request);type(text)"`
	InstanceName string `json:"instance_name"orm:"-"`
}

func (t *Request) TableName() string {
	return "request"
}

func AddRequest(m *Request) (int64, error) {
	o := orm.NewOrm()
	return o.Insert(m)
}

func GetBuildJobCreationRequestByName(name string) (*Request, error) {
	sqlStr := `select * from request where name = ? and type = ? `
	o := orm.NewOrm()
	r := &Request{}
	err := o.Raw(sqlStr, name, common.BuildJobCreateRequestType).QueryRow(&r)
	if err != nil {
		return nil, err
	}
	return r, nil
}