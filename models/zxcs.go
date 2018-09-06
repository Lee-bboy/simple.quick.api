package models

import (
	"fmt"

	"datacenter.analysis.api/common"
	"github.com/go-xorm/xorm"
)

var engine *xorm.Engine

func init() {
	engine = common.RegisterMysql()
}

func GetProductAssociate(num string) []map[string][]byte {
	sql := "select * from zxcs_order_" + num + " limit 4000"
	results, _ := engine.Query(sql)
	fmt.Println(sql)
	//fmt.Println(results)
	return results
}
