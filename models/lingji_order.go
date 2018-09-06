package models

type LingjiOrder struct {
	Id           int64   `xorm:"pk auto_incr"`
	OrderId      string  `xorm:"varchar(50)"`
	Account      string  `xorm:"varchar(255)"`
	DeviceId     string  `xorm:"varchar(128)"`
	Phone        string  `xorm:"varchar(20)"`
	ProductName  string  `xorm:"varchar(255)"`
	ProductType  string  `xorm:"varchar(255)"`
	Price        float64 `xorm:"decimal(12,2)"`
	PayWay       int16   `xorm:"smallint(4)"`
	Channel      string  `xorm:"varchar(50)"`
	OrderApp     int
	TerminalType int8
	IpAdress     string `xorm:"varchar(100)"`
	OrderTime    int
	PayTime      int
	CreatedTime  int
}

//计算两个订单之间支付时间相差值
//以秒为单位
func (o *LingjiOrder) PayTimeInterval(oo *LingjiOrder) int {
	return o.PayTime - oo.PayTime
}

//计算两个订单之间支付时间相差值
//以小时为单位
func (o *LingjiOrder) PayTimeIntervalByHours(oo *LingjiOrder) int {
	return (o.PayTime - oo.PayTime) / 3600
}

type LingjiOrderUser struct {
	LingjiOrder `xorm:"extends"`
	LingjiUser  `xorm:"extends"`
}
