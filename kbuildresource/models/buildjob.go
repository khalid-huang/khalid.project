package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/sirupsen/logrus"
	"time"
)

type Pod struct {
	ID int `orm:"column(id);auto;" description:"只读，主键字段，由后台数据库自动生成"`
	Name string `orm:"column(name);size(256);" description:"等于slave-name"`
	ClusterName string `orm:"column(cluster_name);size(256)" description:"必选，集群名"`
	Labels string `orm:"column(labels);size(256)" description:"标签"`
	Namespace string `orm:"column(namespace);size(256)" description:"命名空间"`
	Status string `orm:"column(status);size(256)" description:"工作负载状态"`
	NodeIP string `orm:"column(node_ip);size(256)" description:"节点ip"`
	IsDelete string `orm:"column(is_delete);default(0)" description:"逻辑删除"`
	GmtCreated time.Time `orm:"column(gmt_created);type(timestamp);auto_now_add;" description:"创建时间"`
	GmtModified time.Time `orm:"column(gmt_modified);type(timestamp);auto_now;" description:"更新更新"`
	Message string `orm:"column(message);" description:"状态运行信息，比如出错原因等，一般是最后一条事件信息"`
	Containers []*Container `orm:"reverse(many)",json:"containers" description:"绑定的containers"`
}

type Container struct {
	ID int `orm:"column(id);auto;" description:"只读，主键字段，由后台数据库自动生成"`
	Pod *Pod `orm:"rel(fk)" json:"pod" description:"只读绑定的podID"`
	CMDs string `json:"cmd" description:"容器启动命令：eg: cm1,cmd2"`
	Name string `json:"name" description:"容器名"`
	Image string `json:"image" description:"容器镜像"`
	RequestCPU string`json:"requestCPU" description:"请求cpu大小"`
	RequestMem string`json:"requestMem" description:"请求内存大小"`
	LimitCPU string`json:"requestCPU" description:"最大可用cpu大小"`
	LimitMem string`json:"requestCPU" description:"最大可用内存大小"`
	GmtCreated time.Time `orm:"column(gmt_created);type(timestamp);auto_now_add;" description:"创建时间"`
	GmtModified time.Time `orm:"column(gmt_modified);type(timestamp);auto_now;" description:"更新更新"`
}

func (t *Pod) TableName() string {
	return "pod"
}

func (t *Container) TableName() string {
	return "container"
}

func AddPod(m *Pod) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	for _, v := range m.Containers {
		var container Container
		container = *v
		container.Pod = m
		_, err = AddPodContainer(&container)
		if err != nil {
			logrus.Info("INFO: insert containers err")
			return
		}
	}
	return
}

func GetPodByID(id int) (v *Pod, err error) {
	o := orm.NewOrm()
	v = &Pod{ID:id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	if err == orm.ErrNoRows {
		return nil,nil
	}
	return nil, err
}

func GetPodByName(name string) (v *Pod, err error) {
	o := orm.NewOrm()
	SQLStr := `select * from pod where name= ? `
	err = o.Raw(SQLStr, name).QueryRow(&v)
	if err = o.Read(v); err == nil {
		return v, nil
	}
	if err == orm.ErrNoRows {
		return nil,nil
	}
	return nil, err
}

// Container
func AddPodContainer(c *Container) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(c)
	return
}