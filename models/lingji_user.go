package models

type LingjiUser struct {
	Id           int64 `xorm:"pk auto_incr"`
	UserId       int64
	Account      string `xorm:"varchar(255)"`
	DeviceId     string `xorm:"varchar(128)"`
	Username     string `xorm:"varchar(255)"`
	Email        string `xorm:"varchar(80)"`
	Phone        string `xorm:"varchar(20)"`
	Wechat       string `xorm:"varchar(50)"`
	Qq           string `xorm:"varchar(20)"`
	Gender       int8
	Birthday     int
	IsWork       int8
	IsMarried    int8
	Province     string `xorm:"varchar(25)"`
	City         string `xorm:"varchar(25)"`
	Region       string `xorm:"varchar(25)"`
	IpAddress    string `xorm:"varchar(50)"`
	BusinessId   int
	RegisterTime int
	CreatedTime  int
}
