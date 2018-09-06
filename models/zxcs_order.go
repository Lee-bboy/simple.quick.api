package models

type ZxcsOrder struct {
	Id          int64   `xorm:"pk autoincr"`
	OrderId     string  `xorm:"varchar(50)"`
	Account     string  `xorm:"varchar(255)"`
	DeviceId    string  `xorm:"varchar(128)"`
	Phone       string  `xorm:"varchar(20)"`
	ProductName string  `xorm:"varchar(255)"`
	ProductType string  `xorm:"varchar(255)"`
	Price       float64 `xorm:"decimal(12,2)"`
	PayWay      string  `xorm:"varchar(50)"`
	Channel     string  `xorm:"varchar(50)"`
	Email       string  `xorm:"varchar(200)"`
	Username    string  `xorm:"varchar(255)"`
	Gender      int8    `xorm:"tinyint(4)"`
	Birthday    int
	OrderApp    int
	IpAddress   string `xorm:"varchar(100)"`
	OrderTime   int
	PayTime     int
	CreatedTime int
}
