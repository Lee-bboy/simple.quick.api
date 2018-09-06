package helper

import (
	"container/list"
	"strconv"
	"time"

	"datacenter.analysis.api/common"
	"datacenter.analysis.api/models"
)

//产品相关性统计，取最前面多少条
const MaxProductCount = 20

type FrequencyResult struct {
	One       int `json:"1"`
	Two       int `json:"2"`
	Three     int `json:"3"`
	Four      int `json:"4"`
	Five      int `json:"5"`
	Six       int `json:"6"`
	Seven     int `json:"7"`
	Eight     int `json:"8"`
	Nine      int `json:"9"`
	Ten       int `json:"10"`
	Eleven    int `json:"11"`
	Twelve    int `json:"12"`
	Thirteen  int `json:"13"`
	Fourteen  int `json:"14"`
	Fifteen   int `json:"15"`
	Sixteen   int `json:"16"`
	Seventeen int `json:"17"`
	Eighteen  int `json:"18"`
	Nineteen  int `json:"19"`
	More      int `json:"20"`
}

func (r *FrequencyResult) Add(rr *FrequencyResult) {
	r.One += rr.One
	r.Two += rr.Two
	r.Three += rr.Three
	r.Four += rr.Four
	r.Five += rr.Five
	r.Six += rr.Six
	r.Seven += rr.Seven
	r.Eight += rr.Eight
	r.Nine += rr.Nine
	r.Ten += rr.Ten
	r.Eleven += rr.Eleven
	r.Twelve += rr.Twelve
	r.Thirteen += rr.Thirteen
	r.Fourteen += rr.Fourteen
	r.Fifteen += rr.Fifteen
	r.Sixteen += rr.Sixteen
	r.Seventeen += rr.Seventeen
	r.Eighteen += rr.Eighteen
	r.Nineteen += rr.Nineteen
	r.More += rr.More
}

type TimeConsume struct {
	One    int `json:"4小时"`
	Two    int `json:"8小时"`
	Three  int `json:"12小时"`
	Four   int `json:"16小时"`
	Five   int `json:"20小时"`
	Six    int `json:"24小时"`
	Seven  int `json:"1-2天"`
	Eight  int `json:"2-3天"`
	Nine   int `json:"3-6天"`
	Ten    int `json:"6-10天"`
	Eleven int `json:"10-14天"`
	Twelve int `json:"14-20天"`
	More   int `json:"20天以上"`
}

func (r *TimeConsume) Add(rr *TimeConsume) {
	r.One += rr.One
	r.Two += rr.Two
	r.Three += rr.Three
	r.Four += rr.Four
	r.Five += rr.Five
	r.Six += rr.Six
	r.Seven += rr.Seven
	r.Eight += rr.Eight
	r.Nine += rr.Nine
	r.Ten += rr.Ten
	r.Eleven += rr.Eleven
	r.Twelve += rr.Twelve
	r.More += rr.More
}

type ProductConsume struct {
	product string
	consume int
}

type AssociateResult map[string]int

func (r AssociateResult) Add(rr AssociateResult) {
	for k, v := range rr {
		r[k] += v
	}
}

type FirstProduct map[string]int

func (r FirstProduct) Add(rr FirstProduct) {
	for k, v := range rr {
		r[k] += v
	}
}

type FirstPrice struct {
	One   int `json:"0-50元"`
	Two   int `json:"51-100元"`
	Three int `json:"101-150元"`
	Four  int `json:"151-200元"`
	Five  int `json:"201-300元"`
	Six   int `json:"301-500元"`
	More  int `json:"500以上"`
}

func (r *FirstPrice) Add(rr *FirstPrice) {
	r.One += rr.One
	r.Two += rr.Two
	r.Three += rr.Three
	r.Four += rr.Four
	r.Five += rr.Five
	r.Six += rr.Six
	r.More += rr.More
}

type FirstTradeResult struct {
	Fproduct FirstProduct `json:"product"`
	Fprice   *FirstPrice  `json:"price"`
}

func LingjiFrequency(data map[string]interface{}) (*FrequencyResult, error) {
	//表名
	tableName := "lingji_order"

	//渠道
	var allChannels bool
	var channels []interface{}
	cc, ok := data["channels"]
	if !ok {
		allChannels = true
	} else {
		allChannels = false
		channels, _ = cc.([]interface{})
	}

	//产品
	var allProducts bool
	var products []interface{}
	pp, ok := data["products"]
	if !ok {
		allProducts = true
	} else {
		allProducts = false
		products, _ = pp.([]interface{})
	}

	//终端类型
	var allTerminalTypes bool
	var terminalTypes []interface{}
	tt, ok := data["terminal_types"]
	if !ok {
		allTerminalTypes = true
	} else {
		allTerminalTypes = false
		terminalTypes, _ = tt.([]interface{})
	}

	//起止时间戳
	var startDate, endDate string
	now := time.Now()
	sd, ok := data["start_date"]
	if !ok {
		startDate = time.Date(now.Year(), now.Month()-1, now.Day(), 0, 0, 0, 0, common.Location).Format("2006-01-02")
	} else {
		startDate, _ = sd.(string)
	}
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", startDate+" 00:00:00", common.Location)
	startTimestamp := startTime.Unix()
	ed, ok := data["end_date"]
	if !ok {
		endDate = now.Format("2006-01-02")
	} else {
		endDate, _ = ed.(string)
	}
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", endDate+" 23:59:59", common.Location)
	endTimestamp := endTime.Unix()

	//获取数据库驱动
	engine := common.RegisterMysql()
	var totalFrequencyResult = new(FrequencyResult)
	var ch = make(chan *struct {
		frequencyResult *FrequencyResult
		err             error
	}, 8)

	//开8个协程处理
	work := func(index int) {
		frequencyResult := new(FrequencyResult)
		result := &struct {
			frequencyResult *FrequencyResult
			err             error
		}{frequencyResult, nil}

		orders := make([]models.LingjiOrder, 0)
		session := engine.Table(tableName+"_"+strconv.Itoa(index)).Where("order_time between ? and ? and pay_time > 0", startTimestamp, endTimestamp)
		session.In("order_app", []int{2000, 438626245})
		if !allChannels {
			session.In("channel", channels)
		}
		if !allProducts {
			session.In("product_name", products)
		}
		if !allTerminalTypes {
			session.In("terminal_type", terminalTypes)
		}
		result.err = session.Cols("account", "device_id").Find(&orders)
		if result.err == nil {
			var tmp = make(map[string]int)

			for _, order := range orders {
				//获取用户唯一表示
				var identifier string
				if order.Account != "" {
					identifier = order.Account
				} else if order.DeviceId != "" {
					identifier = order.DeviceId
				} else {
					continue
				}
				tmp[identifier]++
			}

			for _, v := range tmp {
				if v == 1 {
					frequencyResult.One++
				} else if v == 2 {
					frequencyResult.Two++
				} else if v == 3 {
					frequencyResult.Three++
				} else if v == 4 {
					frequencyResult.Four++
				} else if v == 5 {
					frequencyResult.Five++
				} else if v == 6 {
					frequencyResult.Six++
				} else if v == 7 {
					frequencyResult.Seven++
				} else if v == 8 {
					frequencyResult.Eight++
				} else if v == 9 {
					frequencyResult.Nine++
				} else if v == 10 {
					frequencyResult.Ten++
				} else if v == 11 {
					frequencyResult.Eleven++
				} else if v == 12 {
					frequencyResult.Twelve++
				} else if v == 13 {
					frequencyResult.Thirteen++
				} else if v == 14 {
					frequencyResult.Fourteen++
				} else if v == 15 {
					frequencyResult.Fifteen++
				} else if v == 16 {
					frequencyResult.Sixteen++
				} else if v == 17 {
					frequencyResult.Seventeen++
				} else if v == 18 {
					frequencyResult.Eighteen++
				} else if v == 19 {
					frequencyResult.Nineteen++
				} else if v >= 20 {
					frequencyResult.More++
				}
			}
		}

		ch <- result
	}

	for i := 1; i <= 8; i++ {
		go work(i)
	}

	for j := 1; j <= 8; j++ {
		tmpResult := <-ch
		if tmpResult.err != nil {
			return nil, tmpResult.err
		}

		totalFrequencyResult.Add(tmpResult.frequencyResult)
	}

	return totalFrequencyResult, nil
}

func ZxcsFrequency(data map[string]interface{}) interface{} {
	//表名
	tableName := "zxcs_order"
	//渠道
	var allChannels bool
	var channels []interface{}
	cc, ok := data["channels"]
	if !ok {
		allChannels = true
	} else {
		allChannels = false
		channels, _ = cc.([]interface{})
	}

	//产品
	var allProducts bool
	var products []interface{}
	pp, ok := data["products"]
	if !ok {
		allProducts = true
	} else {
		allProducts = false
		products, _ = pp.([]interface{})
	}

	//起止时间戳
	var startDate, endDate string
	now := time.Now()
	sd, ok := data["start_date"]
	if !ok {
		startDate = time.Date(now.Year(), now.Month()-1, now.Day(), 0, 0, 0, 0, common.Location).Format("2006-01-02")
	} else {
		startDate, _ = sd.(string)
	}
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", startDate+" 00:00:00", common.Location)
	startTimestamp := startTime.Unix()
	ed, ok := data["end_date"]
	if !ok {
		endDate = now.Format("2006-01-02")
	} else {
		endDate, _ = ed.(string)
	}

	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", endDate+" 23:59:59", common.Location)
	endTimestamp := endTime.Unix()

	//获取数据库驱动
	engine := common.RegisterMysql()

	var ch = make(chan map[string]int)

	for i := 1; i <= 8; i++ {
		go func(index int) {
			users := make(map[string]int)
			orders := make([]models.ZxcsOrder, 0)
			session := engine.Table(tableName+"_"+strconv.Itoa(index)).Where("order_time between ? and ? AND pay_time > 0", startTimestamp, endTimestamp)
			if !allChannels {
				session.In("channel", channels)
			}
			if !allProducts {
				session.In("product_name", products)
			}

			err := session.Cols("email").Find(&orders)

			if err == nil && len(orders) > 0 {
				for _, order := range orders {
					if _, ok := users[order.Email]; ok {
						users[order.Email] = users[order.Email] + 1
					} else {
						users[order.Email] = 1
					}
				}
			}

			ch <- users
		}(i)
	}

	result := new(FrequencyResult)

	for j := 1; j <= 8; j++ {
		tmpResult := <-ch

		if len(tmpResult) > 0 {
			for _, value := range tmpResult {

				if value == 1 {
					result.One++
				} else if value == 2 {
					result.Two++
				} else if value == 3 {
					result.Three++
				} else if value == 4 {
					result.Four++
				} else if value == 5 {
					result.Five++
				} else if value == 6 {
					result.Six++
				} else if value == 7 {
					result.Seven++
				} else if value == 8 {
					result.Eight++
				} else if value == 9 {
					result.Nine++
				} else if value == 10 {
					result.Ten++
				} else if value == 11 {
					result.Eleven++
				} else if value == 12 {
					result.Twelve++
				} else if value == 13 {
					result.Thirteen++
				} else if value == 14 {
					result.Fourteen++
				} else if value == 15 {
					result.Fifteen++
				} else if value == 16 {
					result.Sixteen++
				} else if value == 17 {
					result.Seventeen++
				} else if value == 18 {
					result.Eighteen++
				} else if value == 19 {
					result.Nineteen++
				} else if value >= 20 {
					result.More++
				}
			}
		}
	}

	return result
}

func ShopFrequency(data map[string]interface{}) {

}

func LingjiTimeInterval(data map[string]interface{}) (*TimeConsume, error) {
	//表名
	tableName := "lingji_order"

	//产品
	var allProducts bool
	var products []interface{}
	pp, ok := data["products"]
	if !ok {
		allProducts = true
	} else {
		allProducts = false
		products, _ = pp.([]interface{})
	}

	//终端类型
	var allTerminalTypes bool
	var terminalTypes []interface{}
	tt, ok := data["terminal_types"]
	if !ok {
		allTerminalTypes = true
	} else {
		allTerminalTypes = false
		terminalTypes, _ = tt.([]interface{})
	}

	//起止时间戳
	var startDate, endDate string
	now := time.Now()
	sd, ok := data["start_date"]
	if !ok {
		startDate = time.Date(now.Year(), now.Month()-1, now.Day(), 0, 0, 0, 0, common.Location).Format("2006-01-02")
	} else {
		startDate, _ = sd.(string)
	}
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", startDate+" 00:00:00", common.Location)
	startTimestamp := startTime.Unix()
	ed, ok := data["end_date"]
	if !ok {
		endDate = now.Format("2006-01-02")
	} else {
		endDate, _ = ed.(string)
	}
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", endDate+" 23:59:59", common.Location)
	endTimestamp := endTime.Unix()

	//获取数据库驱动
	engine := common.RegisterMysql()
	var totalTimeIntervalResult = new(TimeConsume)
	var ch = make(chan *struct {
		timeIntervalResult *TimeConsume
		err                error
	}, 8)

	//开8个协程处理
	work := func(index int) {
		var timeIntervalResult = new(TimeConsume)
		result := &struct {
			timeIntervalResult *TimeConsume
			err                error
		}{timeIntervalResult, nil}

		orders := make([]*models.LingjiOrder, 0)
		session := engine.Table(tableName+"_"+strconv.Itoa(index)).Where("order_time between ? and ? and pay_time > 0", startTimestamp, endTimestamp)
		session.In("order_app", []int{2000, 438626245})
		if !allProducts {
			session.In("product_name", products)
		}
		if !allTerminalTypes {
			session.In("terminal_type", terminalTypes)
		}
		result.err = session.Cols("account", "device_id", "pay_time").Asc("pay_time").Find(&orders)
		if result.err == nil {
			var user = make(map[string][]*models.LingjiOrder)

			for _, order := range orders {
				//获取用户唯一表示
				var identifier string
				if order.Account != "" {
					identifier = order.Account
				} else if order.DeviceId != "" {
					identifier = order.DeviceId
				} else {
					continue
				}

				if _, ok := user[identifier]; !ok {
					user[identifier] = make([]*models.LingjiOrder, 0)
				}
				if len(user[identifier]) < 2 {
					user[identifier] = append(user[identifier], order)
				}
			}

			for _, v := range user {
				//跳过消费不到两次的用户
				if len(v) < 2 {
					continue
				}

				//计算用户消费时间间隔
				time := v[1].PayTimeIntervalByHours(v[0])
				switch {
				case 0 <= time && time <= 4:
					timeIntervalResult.One++
				case 4 < time && time <= 8:
					timeIntervalResult.Two++
				case 8 < time && time <= 12:
					timeIntervalResult.Three++
				case 12 < time && time <= 16:
					timeIntervalResult.Four++
				case 16 < time && time <= 20:
					timeIntervalResult.Five++
				case 20 < time && time <= 24:
					timeIntervalResult.Six++
				case 24 < time && time <= 48:
					timeIntervalResult.Seven++
				case 48 < time && time <= 72:
					timeIntervalResult.Eight++
				case 72 < time && time <= 144:
					timeIntervalResult.Nine++
				case 144 < time && time <= 240:
					timeIntervalResult.Ten++
				case 240 < time && time <= 336:
					timeIntervalResult.Eleven++
				case 336 < time && time <= 480:
					timeIntervalResult.Twelve++
				case 480 < time:
					timeIntervalResult.More++
				}
			}
		}

		ch <- result
	}

	for i := 1; i <= 8; i++ {
		go work(i)
	}

	for j := 1; j <= 8; j++ {
		tmpResult := <-ch
		if tmpResult.err != nil {
			return nil, tmpResult.err
		}

		totalTimeIntervalResult.Add(tmpResult.timeIntervalResult)
	}

	return totalTimeIntervalResult, nil
}

func ZxcsTimeInterval(data map[string]interface{}) interface{} {
	//表名
	tableName := "zxcs_order"

	//产品
	var allProducts bool
	var products []interface{}
	pp, ok := data["products"]
	if !ok {
		allProducts = true
	} else {
		allProducts = false
		products, _ = pp.([]interface{})
	}

	//起止时间戳
	var startDate, endDate string
	now := time.Now()
	sd, ok := data["start_date"]
	if !ok {
		startDate = time.Date(now.Year(), now.Month()-1, now.Day(), 0, 0, 0, 0, common.Location).Format("2006-01-02")
	} else {
		startDate, _ = sd.(string)
	}
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", startDate+" 00:00:00", common.Location)
	startTimestamp := startTime.Unix()
	ed, ok := data["end_date"]
	if !ok {
		endDate = now.Format("2006-01-02")
	} else {
		endDate, _ = ed.(string)
	}

	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", endDate+" 23:59:59", common.Location)
	endTimestamp := endTime.Unix()

	is_two := 1 //第二次付费间隔-1/平均消费间隔-2

	//获取数据库驱动
	engine := common.RegisterMysql()

	var ch = make(chan map[string]map[string]int)

	work := func(index int) {

		users := make(map[string]map[string]int)

		orders := make([]models.ZxcsOrder, 0)
		session := engine.Table(tableName+"_"+strconv.Itoa(index)).Where("order_time between ? and ? AND pay_time > 0", startTimestamp, endTimestamp)

		if !allProducts {
			session.In("product_name", products)
		}

		err := session.Cols("email", "pay_time").Asc("pay_time").Find(&orders)

		if err == nil && len(orders) > 0 {
			for _, order := range orders {
				if t, ok := users[order.Email]; ok {

					if is_two == 1 {
						if t["num"] >= 2 {
							continue
						}
						users[order.Email]["num"] = t["num"] + 1
						users[order.Email]["time"] = order.PayTime - t["time"]
					} else {
						users[order.Email]["num"] = t["num"] + 1
						users[order.Email]["time"] = order.PayTime + t["time"]
					}

				} else {
					consume := make(map[string]int)
					consume["num"] = 1
					consume["time"] = order.PayTime
					users[order.Email] = consume
				}
			}
		}

		ch <- users
	}

	for i := 1; i <= 8; i++ {
		go work(i)
	}

	result := new(TimeConsume)

	for j := 1; j <= 8; j++ {
		tmpResult := <-ch

		if is_two == 1 {
			for _, value := range tmpResult {

				if len(value) > 0 {
					if value["num"] > 1 {
						time := value["time"] / 3600

						//判断时间间隔
						switch {
						case 0 <= time && time <= 4:
							result.One++
						case 4 < time && time <= 8:
							result.Two++
						case 8 < time && time <= 12:
							result.Three++
						case 12 < time && time <= 16:
							result.Four++
						case 16 < time && time <= 20:
							result.Five++
						case 20 < time && time <= 24:
							result.Six++
						case 24 < time && time <= 48: //1-2
							result.Seven++
						case 48 < time && time <= 72: //2-3
							result.Eight++
						case 72 < time && time <= 144: //3-6
							result.Nine++
						case 144 < time && time <= 240: //6-10
							result.Ten++
						case 240 < time && time <= 336: //10-14
							result.Eleven++
						case 336 < time && time <= 480: //14-20
							result.Twelve++
						case 480 < time:
							result.More++
						}
					}
				}
			}
		} else {
			//todo
		}
	}

	return result
}

func LingjiAssociate(data map[string]interface{}) (AssociateResult, error) {
	//表名
	tableName := "lingji_order"

	//渠道
	var allChannels bool
	var channels []interface{}
	cc, ok := data["channels"]
	if !ok {
		allChannels = true
	} else {
		allChannels = false
		channels, _ = cc.([]interface{})
	}

	//产品
	var products []interface{}
	pp, _ := data["products"]
	products, _ = pp.([]interface{})

	//终端类型
	var allTerminalTypes bool
	var terminalTypes []interface{}
	tt, ok := data["terminal_types"]
	if !ok {
		allTerminalTypes = true
	} else {
		allTerminalTypes = false
		terminalTypes, _ = tt.([]interface{})
	}

	//起止时间戳
	var startDate, endDate string
	now := time.Now()
	sd, ok := data["start_date"]
	if !ok {
		startDate = time.Date(now.Year(), now.Month()-1, now.Day(), 0, 0, 0, 0, common.Location).Format("2006-01-02")
	} else {
		startDate, _ = sd.(string)
	}
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", startDate+" 00:00:00", common.Location)
	startTimestamp := startTime.Unix()
	ed, ok := data["end_date"]
	if !ok {
		endDate = now.Format("2006-01-02")
	} else {
		endDate, _ = ed.(string)
	}
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", endDate+" 23:59:59", common.Location)
	endTimestamp := endTime.Unix()

	//获取数据库驱动
	engine := common.RegisterMysql()
	var totalAssociateResult AssociateResult = make(map[string]int)
	var ch = make(chan *struct {
		associateResult AssociateResult
		err             error
	}, 8)

	//开8个协程处理
	work := func(index int) {
		var associateResult AssociateResult = make(map[string]int)
		result := &struct {
			associateResult AssociateResult
			err             error
		}{associateResult, nil}

		orders := make([]models.LingjiOrder, 0)
		session := engine.Table(tableName+"_"+strconv.Itoa(index)).Where("order_time between ? and ? and pay_time > 0", startTimestamp, endTimestamp)
		session.In("order_app", []int{2000, 438626245})
		session.In("product_name", products)
		if !allChannels {
			session.In("channel", channels)
		}
		if !allTerminalTypes {
			session.In("terminal_type", terminalTypes)
		}
		result.err = session.Cols("account", "device_id").Find(&orders)
		if result.err == nil {
			var accounts = make([]string, 0)
			var accountAlready = make(map[string]bool)
			var deviceIds = make([]string, 0)
			var deviceIdAlready = make(map[string]bool)

			for _, order := range orders {
				if order.Account != "" {
					if _, ok1 := accountAlready[order.Account]; !ok1 {
						accountAlready[order.Account] = true
						accounts = append(accounts, order.Account)
					}
				} else if order.DeviceId != "" {
					if _, ok2 := deviceIdAlready[order.DeviceId]; !ok2 {
						deviceIdAlready[order.DeviceId] = true
						deviceIds = append(deviceIds, order.DeviceId)
					}
				} else {
					continue
				}
			}

			var productAlready map[string]map[string]bool
			if len(accounts) > 0 {
				orders = make([]models.LingjiOrder, 0)
				session := engine.Table(tableName+"_"+strconv.Itoa(index)).Where("order_time between ? and ? and pay_time > 0", startTimestamp, endTimestamp)
				session.In("order_app", []int{2000, 438626245})
				session.NotIn("product_name", products)
				session.In("account", accounts)
				if !allChannels {
					session.In("channel", channels)
				}
				if !allTerminalTypes {
					session.In("terminal_type", terminalTypes)
				}
				result.err = session.Cols("account", "product_name").Find(&orders)
				if result.err == nil {
					productAlready = make(map[string]map[string]bool)
					for _, order := range orders {
						if _, ok1 := productAlready[order.ProductName]; !ok1 {
							productAlready[order.ProductName] = make(map[string]bool)
						}
						if _, ok2 := productAlready[order.ProductName][order.Account]; !ok2 {
							productAlready[order.ProductName][order.Account] = true
							associateResult[order.ProductName]++
						}
					}
				}
			}

			if len(deviceIds) > 0 && result.err == nil {
				orders = make([]models.LingjiOrder, 0)
				session := engine.Table(tableName+"_"+strconv.Itoa(index)).Where("order_time between ? and ? and pay_time > 0", startTimestamp, endTimestamp)
				session.In("order_app", []int{2000, 438626245})
				session.NotIn("product_name", products)
				session.In("device_id", deviceIds)
				if !allChannels {
					session.In("channel", channels)
				}
				if !allTerminalTypes {
					session.In("terminal_type", terminalTypes)
				}
				result.err = session.Cols("device_id", "product_name").Find(&orders)
				if result.err == nil {
					productAlready = make(map[string]map[string]bool)
					for _, order := range orders {
						if _, ok1 := productAlready[order.ProductName]; !ok1 {
							productAlready[order.ProductName] = make(map[string]bool)
						}
						if _, ok2 := productAlready[order.ProductName][order.DeviceId]; !ok2 {
							productAlready[order.ProductName][order.DeviceId] = true
							associateResult[order.ProductName]++
						}
					}
				}
			}
		}

		ch <- result
	}

	for i := 1; i <= 8; i++ {
		go work(i)
	}

	for j := 1; j <= 8; j++ {
		tmpResult := <-ch
		if tmpResult.err != nil {
			return nil, tmpResult.err
		}

		totalAssociateResult.Add(tmpResult.associateResult)
	}

	//按消费用户数降序排列
	l := list.New()
	for k, v := range totalAssociateResult {
		if l.Len() == 0 {
			l.PushFront(&ProductConsume{k, v})
		} else {
			smallest := true
			for node := l.Front(); node != nil; node = node.Next() {
				value := node.Value.(*ProductConsume)
				if v >= value.consume {
					l.InsertBefore(&ProductConsume{k, v}, node)
					smallest = false
					break
				}
			}

			if smallest {
				l.PushBack(&ProductConsume{k, v})
			}

			//如果列表超过最大长度，移除最后面的元素
			if l.Len() > MaxProductCount {
				l.Remove(l.Back())
			}
		}
	}

	var topAssociateResult AssociateResult = make(map[string]int)
	for node := l.Front(); node != nil; node = node.Next() {
		value := node.Value.(*ProductConsume)
		topAssociateResult[value.product] = value.consume
	}

	return topAssociateResult, nil
}

func ZxcsAssociate(data map[string]interface{}) interface{} {
	//表名
	tableName := "zxcs_order"

	//产品
	var products []interface{}
	pp, ok := data["products"]

	if ok {
		//默认八字精批
		products, _ = pp.([]interface{})
	}

	//渠道
	var allChannels bool
	var channels []interface{}
	cc, ok := data["channels"]
	if !ok {
		allChannels = true
	} else {
		allChannels = false
		channels, _ = cc.([]interface{})
	}

	//起止时间戳
	var startDate, endDate string
	now := time.Now()
	sd, ok := data["start_date"]
	if !ok {
		startDate = time.Date(now.Year(), now.Month()-1, now.Day(), 0, 0, 0, 0, common.Location).Format("2006-01-02")
	} else {
		startDate, _ = sd.(string)
	}

	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", startDate+" 00:00:00", common.Location)
	startTimestamp := startTime.Unix()
	ed, ok := data["end_date"]
	if !ok {
		endDate = now.Format("2006-01-02")
	} else {
		endDate, _ = ed.(string)
	}

	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", endDate+" 23:59:59", common.Location)
	endTimestamp := endTime.Unix()

	//获取数据库驱动
	engine := common.RegisterMysql()

	var ch = make(chan map[int][]string)
	work := func(index int) {

		e := make(map[int][]string)
		users := make(map[string]bool)

		orders := make([]models.ZxcsOrder, 0)
		session := engine.Table(tableName+"_"+strconv.Itoa(index)).Where("order_time between ? and ? AND pay_time > 0", startTimestamp, endTimestamp)

		session.In("product_name", products)

		if !allChannels {
			session.In("channel", channels)
		}

		err := session.Cols("email").Find(&orders)

		if err == nil && len(orders) > 0 {
			for _, order := range orders {
				if _, ok := users[order.Email]; ok {
					continue
				} else {
					users[order.Email] = true
					e[index] = append(e[index], order.Email)
				}
			}
		}

		ch <- e
	}

	for i := 1; i <= 8; i++ {
		go work(i)
	}

	var ch1 = make(chan map[string]map[string]bool)

	work_p := func(index int) {
		tmpEmail := <-ch

		orders := make([]models.ZxcsOrder, 0)

		productMap := make(map[string]map[string]bool)

		for key, emails := range tmpEmail {
			if len(emails) > 0 {

				session := engine.Table(tableName+"_"+strconv.Itoa(key)).Where("order_time between ? and ? AND pay_time > 0", startTimestamp, endTimestamp)

				session.In("email", emails)

				if !allChannels {
					session.In("channel", channels)
				}

				session.NotIn("product_name", products)

				err := session.Cols("email", "product_name").Find(&orders)

				if err == nil {
					for _, order := range orders {
						if p, ok := productMap[order.ProductName]; ok {
							if _, ok := p[order.Email]; ok {
								continue
							} else {
								p[order.Email] = true
								productMap[order.ProductName] = p
							}
						} else {
							emailMap := make(map[string]bool)
							emailMap[order.Email] = true
							productMap[order.ProductName] = emailMap
						}
					}
				}
			}
		}

		ch1 <- productMap
	}

	for j := 1; j <= 8; j++ {
		go work_p(j)
	}

	result := make(map[string]int)

	for k := 1; k <= 8; k++ {
		tmpProduct := <-ch1

		for index, value := range tmpProduct {
			if len(value) < 1 {
				continue
			}

			if r, ok := result[index]; ok {
				result[index] = r + len(value)
			} else {
				result[index] = len(value)
			}
		}
	}

	if len(result) < 1 {
		return nil
	}

	return result

}

func LingjiFirstTrade(data map[string]interface{}) (*FirstTradeResult, error) {
	//表名
	tableName := "lingji_order"

	//终端类型
	var allTerminalTypes bool
	var terminalTypes []interface{}
	tt, ok := data["terminal_types"]
	if !ok {
		allTerminalTypes = true
	} else {
		allTerminalTypes = false
		terminalTypes, _ = tt.([]interface{})
	}

	//起止时间戳
	var startDate, endDate string
	now := time.Now()
	sd, ok := data["start_date"]
	if !ok {
		startDate = time.Date(now.Year(), now.Month()-1, now.Day(), 0, 0, 0, 0, common.Location).Format("2006-01-02")
	} else {
		startDate, _ = sd.(string)
	}
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", startDate+" 00:00:00", common.Location)
	startTimestamp := startTime.Unix()
	ed, ok := data["end_date"]
	if !ok {
		endDate = now.Format("2006-01-02")
	} else {
		endDate, _ = ed.(string)
	}
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", endDate+" 23:59:59", common.Location)
	endTimestamp := endTime.Unix()

	//获取数据库驱动
	engine := common.RegisterMysql()
	var firstTradeResult = new(FirstTradeResult)
	firstTradeResult.Fproduct = make(map[string]int)
	firstTradeResult.Fprice = new(FirstPrice)
	var ch = make(chan *struct {
		firstProduct FirstProduct
		firstPrice   *FirstPrice
		err          error
	}, 8)

	//开8个协程处理
	work := func(index int) {
		var firstProduct = make(map[string]int)
		firstPrice := new(FirstPrice)
		result := &struct {
			firstProduct FirstProduct
			firstPrice   *FirstPrice
			err          error
		}{firstProduct, firstPrice, nil}

		orders := make([]models.LingjiOrder, 0)
		session := engine.Table(tableName+"_"+strconv.Itoa(index)).Where("order_time between ? and ? and pay_time > 0", startTimestamp, endTimestamp)
		session.In("order_app", []int{2000, 438626245})
		if !allTerminalTypes {
			session.In("terminal_type", terminalTypes)
		}
		result.err = session.Cols("account", "device_id", "product_name", "price").Asc("pay_time").Find(&orders)
		if result.err == nil {
			already := make(map[string]bool)

			for _, order := range orders {
				//获取用户唯一表示
				var identifier string
				if order.Account != "" {
					identifier = order.Account
				} else if order.DeviceId != "" {
					identifier = order.DeviceId
				} else {
					continue
				}

				if _, ok := already[identifier]; !ok {
					already[identifier] = true

					firstProduct[order.ProductName]++

					switch {
					case order.Price < 51:
						firstPrice.One++
					case order.Price >= 51 && order.Price < 101:
						firstPrice.Two++
					case order.Price >= 101 && order.Price < 151:
						firstPrice.Three++
					case order.Price >= 151 && order.Price < 201:
						firstPrice.Four++
					case order.Price >= 201 && order.Price < 301:
						firstPrice.Five++
					case order.Price >= 301 && order.Price < 501:
						firstPrice.Six++
					case order.Price >= 501:
						firstPrice.More++
					}
				}
			}
		}

		ch <- result
	}

	for i := 1; i <= 8; i++ {
		go work(i)
	}

	for j := 1; j <= 8; j++ {
		tmpResult := <-ch
		if tmpResult.err != nil {
			return nil, tmpResult.err
		}

		firstTradeResult.Fproduct.Add(tmpResult.firstProduct)
		firstTradeResult.Fprice.Add(tmpResult.firstPrice)
	}

	//按消费用户数降序排列
	l := list.New()
	for k, v := range firstTradeResult.Fproduct {
		if l.Len() == 0 {
			l.PushFront(&ProductConsume{k, v})
		} else {
			smallest := true
			for node := l.Front(); node != nil; node = node.Next() {
				value := node.Value.(*ProductConsume)
				if v >= value.consume {
					l.InsertBefore(&ProductConsume{k, v}, node)
					smallest = false
					break
				}
			}

			if smallest {
				l.PushBack(&ProductConsume{k, v})
			}

			//如果列表超过最大长度，移除最后面的元素
			if l.Len() > MaxProductCount {
				l.Remove(l.Back())
			}
		}
	}

	var topFirstProduct = make(map[string]int)
	for node := l.Front(); node != nil; node = node.Next() {
		value := node.Value.(*ProductConsume)
		topFirstProduct[value.product] = value.consume
	}
	firstTradeResult.Fproduct = topFirstProduct

	return firstTradeResult, nil
}
