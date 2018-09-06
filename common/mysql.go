package common

import (
	"fmt"

	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var engine *xorm.Engine
var mutex = new(sync.RWMutex)

func RegisterMysql() *xorm.Engine {
	mutex.RLock()
	if engine != nil {
		mutex.RUnlock()
		return engine
	} else {
		mutex.RUnlock()
		mutex.Lock()

		if engine == nil {
			var dns string

			db_host := GetConf("db.host", "")  //填入默认配置
			db_port := GetConf("db.port", "")
			db_user := GetConf("db.user", "")
			db_pass := GetConf("db.pass", "")
			db_name := GetConf("db.name", "")

			dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", db_user, db_pass, db_host, db_port, db_name)

			var err error
			engine, err = xorm.NewEngine("mysql", dns)
			if err != nil {
				fmt.Println("连接数据库失败")
			}

			engine.SetMaxIdleConns(10)  //连接池的空闲数大小
			engine.SetMaxOpenConns(100) //最大打开连接数
		}

		mutex.Unlock()
		return engine
	}
}
