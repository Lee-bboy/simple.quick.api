package controllers

import (
	"net/http"
	"strconv"
	"time"

	"datacenter.analysis.api/common"
	"datacenter.analysis.api/models"
)

var UPController = new(UserPortraitController)

//消费人数
type GenderResult struct {
	Male   int `json:"male"`
	Female int `json:"female"`
	Unkown int `json:"unkown"`
	All    int `json:"all"`
}

//消费转化率
type GenderConversion struct {
	Male   float64 `json:"male"`
	Female float64 `json:"female"`
	Unkown float64 `json:"unkown"`
	All    float64 `json:"all"`
}

func (gr *GenderResult) Add(gr2 *GenderResult) {
	gr.Male += gr2.Male
	gr.Female += gr2.Female
	gr.Unkown += gr2.Unkown
	gr.All += gr2.All
}

func calculateGenderConversion(paidGenderResult, allGenderResult *GenderResult) (conversion *GenderConversion) {
	conversion = new(GenderConversion)

	if paidGenderResult.Male == 0 {
		conversion.Male = 0
	} else {
		conversion.Male = float64(paidGenderResult.Male) / float64(allGenderResult.Male)
	}

	if paidGenderResult.Female == 0 {
		conversion.Female = 0
	} else {
		conversion.Female = float64(paidGenderResult.Female) / float64(allGenderResult.Female)
	}

	if paidGenderResult.Unkown == 0 {
		conversion.Unkown = 0
	} else {
		conversion.Unkown = float64(paidGenderResult.Unkown) / float64(allGenderResult.Unkown)
	}

	if paidGenderResult.All == 0 {
		conversion.All = 0
	} else {
		conversion.All = float64(paidGenderResult.All) / float64(allGenderResult.All)
	}

	return conversion
}

type AgeResult map[string]int

type AgeConversion map[string]float64

func (ar AgeResult) Add(ar2 AgeResult) {
	for k, v := range ar2 {
		ar[k] += v
	}
}

func calculateAgeConversion(paidAgeResult, allAgeResult AgeResult) (conversion AgeConversion) {
	conversion = make(map[string]float64)

	for k, v := range allAgeResult {
		paidCount, ok := paidAgeResult[k]
		if !ok || paidCount == 0 {
			conversion[k] = 0
		} else {
			conversion[k] = float64(paidCount) / float64(v)
		}
	}

	return conversion
}

type HourResult map[string]int

type HourConversion map[string]float64

func (hr HourResult) Add(hr2 HourResult) {
	for k, v := range hr2 {
		hr[k] += v
	}
}

func calculateHourConversion(paidHourResult, allHourResult HourResult) (conversion HourConversion) {
	conversion = make(map[string]float64)

	for k, v := range allHourResult {
		paidCount, ok := paidHourResult[k]
		if !ok || paidCount == 0 {
			conversion[k] = 0
		} else {
			conversion[k] = float64(paidCount) / float64(v)
		}
	}

	return conversion
}

type UserPortraitController struct {
	Controller
}

func (c *UserPortraitController) Gender(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	data, err := c.parseParams(r)
	if err != nil {
		c.errorResponse(w, 400, "请求不正确")
		return
	}

	//获取业务线
	bid, ok := data["business_id"]
	if !ok {
		c.errorResponse(w, 400, "请求不正确")
		return
	}
	businessId, ok := bid.(float64)
	if !ok {
		c.errorResponse(w, 400, "请求1不正确")
		return
	}

	//验证products参数
	p, ok := data["products"]
	if ok {
		products, ok := p.([]interface{})
		if !ok {
			c.errorResponse(w, 400, "请求不正确")
			return
		}
		for _, product := range products {
			_, ok := product.(string)
			if !ok {
				c.errorResponse(w, 400, "请求不正确")
				return
			}
		}
	}

	//验证channels参数
	cc, ok := data["channels"]
	if ok {
		channels, ok := cc.([]interface{})
		if !ok {
			c.errorResponse(w, 400, "请求不正确")
			return
		}
		for _, channel := range channels {
			_, ok := channel.(string)
			if !ok {
				c.errorResponse(w, 400, "请求不正确")
				return
			}
		}
	}

	//验证terminal_types参数
	if businessId == 1 {
		tt1, ok := data["terminal_types"]
		if ok {
			tt2, ok := tt1.([]interface{})
			if !ok {
				c.errorResponse(w, 400, "请求不正确")
				return
			}
			for _, tt := range tt2 {
				switch tt3 := tt.(type) {
				case float64:
				case string:
					_, err := strconv.Atoi(tt3)
					if err != nil {
						c.errorResponse(w, 400, "请求不正确")
						return
					}
				default:
					c.errorResponse(w, 400, "请求不正确")
					return
				}
			}
		}
	}

	//验证时间参数
	sd, ok := data["start_date"]
	if ok {
		startDate, ok := sd.(string)
		if !ok {
			c.errorResponse(w, 400, "请求不正确")
			return
		}
		_, err = time.ParseInLocation("2006-01-02", startDate, common.Location)
		if err != nil {
			c.errorResponse(w, 400, "请求不正确")
			return
		}
	}
	ed, ok := data["end_date"]
	if ok {
		endDate, ok := ed.(string)
		if !ok {
			c.errorResponse(w, 400, "请求不正确")
			return
		}
		_, err = time.ParseInLocation("2006-01-02", endDate, common.Location)
		if err != nil {
			c.errorResponse(w, 400, "请求不正确")
			return
		}
	}

	switch businessId {
	case 1:
		c.lingjiGender(w, r, data)
	case 2:
		c.zxcsGender(w, r, data)
	case 3:
		c.shopGender(w, r, data)
	default:
		c.errorResponse(w, 400, "请求不正确")
		return
	}
}

func (c *UserPortraitController) Age(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	data, err := c.parseParams(r)
	if err != nil {
		c.errorResponse(w, 400, "请求不正确")
		return
	}

	//获取业务线
	bid, ok := data["business_id"]
	if !ok {
		c.errorResponse(w, 400, "请求不正确")
		return
	}
	businessId, ok := bid.(float64)
	if !ok {
		c.errorResponse(w, 400, "请求不正确")
		return
	}

	//验证products参数
	p, ok := data["products"]
	if ok {
		products, ok := p.([]interface{})
		if !ok {
			c.errorResponse(w, 400, "请求不正确")
			return
		}
		for _, product := range products {
			_, ok := product.(string)
			if !ok {
				c.errorResponse(w, 400, "请求不正确")
				return
			}
		}
	}

	//验证channels参数
	cc, ok := data["channels"]
	if ok {
		channels, ok := cc.([]interface{})
		if !ok {
			c.errorResponse(w, 400, "请求不正确")
			return
		}
		for _, channel := range channels {
			_, ok := channel.(string)
			if !ok {
				c.errorResponse(w, 400, "请求不正确")
				return
			}
		}
	}

	//验证terminal_types参数
	if businessId == 1 {
		tt1, ok := data["terminal_types"]
		if ok {
			tt2, ok := tt1.([]interface{})
			if !ok {
				c.errorResponse(w, 400, "请求不正确")
				return
			}
			for _, tt := range tt2 {
				switch tt3 := tt.(type) {
				case float64:
				case string:
					_, err := strconv.Atoi(tt3)
					if err != nil {
						c.errorResponse(w, 400, "请求不正确")
						return
					}
				default:
					c.errorResponse(w, 400, "请求不正确")
					return
				}
			}
		}
	}

	//验证时间参数
	sd, ok := data["start_date"]
	if ok {
		startDate, ok := sd.(string)
		if !ok {
			c.errorResponse(w, 400, "请求不正确")
			return
		}
		_, err = time.ParseInLocation("2006-01-02", startDate, common.Location)
		if err != nil {
			c.errorResponse(w, 400, "请求不正确")
			return
		}
	}
	ed, ok := data["end_date"]
	if ok {
		endDate, ok := ed.(string)
		if !ok {
			c.errorResponse(w, 400, "请求不正确")
			return
		}
		_, err = time.ParseInLocation("2006-01-02", endDate, common.Location)
		if err != nil {
			c.errorResponse(w, 400, "请求不正确")
			return
		}
	}

	switch businessId {
	case 1:
		c.lingjiAge(w, r, data)
	case 2:
		c.zxcsAge(w, r, data)
	case 3:
		c.shopAge(w, r, data)
	default:
		c.errorResponse(w, 400, "请求不正确")
		return
	}
}

func (c *UserPortraitController) Hour(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	data, err := c.parseParams(r)
	if err != nil {
		c.errorResponse(w, 400, "请求不正确")
		return
	}

	//获取业务线
	bid, ok := data["business_id"]
	if !ok {
		c.errorResponse(w, 400, "请求不正确")
		return
	}
	businessId, ok := bid.(float64)
	if !ok {
		c.errorResponse(w, 400, "请求不正确")
		return
	}

	//验证products参数
	p, ok := data["products"]
	if ok {
		products, ok := p.([]interface{})
		if !ok {
			c.errorResponse(w, 400, "请求不正确")
			return
		}
		for _, product := range products {
			_, ok := product.(string)
			if !ok {
				c.errorResponse(w, 400, "请求不正确")
				return
			}
		}
	}

	//验证channels参数
	cc, ok := data["channels"]
	if ok {
		channels, ok := cc.([]interface{})
		if !ok {
			c.errorResponse(w, 400, "请求不正确")
			return
		}
		for _, channel := range channels {
			_, ok := channel.(string)
			if !ok {
				c.errorResponse(w, 400, "请求不正确")
				return
			}
		}
	}

	//验证terminal_types参数
	if businessId == 1 {
		tt1, ok := data["terminal_types"]
		if ok {
			tt2, ok := tt1.([]interface{})
			if !ok {
				c.errorResponse(w, 400, "请求不正确")
				return
			}
			for _, tt := range tt2 {
				switch tt3 := tt.(type) {
				case float64:
				case string:
					_, err := strconv.Atoi(tt3)
					if err != nil {
						c.errorResponse(w, 400, "请求不正确")
						return
					}
				default:
					c.errorResponse(w, 400, "请求不正确")
					return
				}
			}
		}
	}

	//验证时间参数
	sd, ok := data["start_date"]
	if ok {
		startDate, ok := sd.(string)
		if !ok {
			c.errorResponse(w, 400, "请求不正确")
			return
		}
		_, err = time.ParseInLocation("2006-01-02", startDate, common.Location)
		if err != nil {
			c.errorResponse(w, 400, "请求不正确")
			return
		}
	}
	ed, ok := data["end_date"]
	if ok {
		endDate, ok := ed.(string)
		if !ok {
			c.errorResponse(w, 400, "请求不正确")
			return
		}
		_, err = time.ParseInLocation("2006-01-02", endDate, common.Location)
		if err != nil {
			c.errorResponse(w, 400, "请求不正确")
			return
		}
	}

	switch businessId {
	case 1:
		c.lingjiHour(w, r, data)
	case 2:
		c.zxcsHour(w, r, data)
	case 3:
		c.shopHour(w, r, data)
	default:
		c.errorResponse(w, 400, "请求不正确")
		return
	}
}

func (c *UserPortraitController) lingjiGender(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
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
	var paidGenderResult = new(GenderResult)
	var allGenderResult = new(GenderResult)
	var ch = make(chan *struct {
		paidGenderResult *GenderResult
		allGenderResult  *GenderResult
		err              error
	}, 8)

	//开8个协程处理
	work := func(index int) {
		paidResult := new(GenderResult)
		allResult := new(GenderResult)
		result := &struct {
			paidGenderResult *GenderResult
			allGenderResult  *GenderResult
			err              error
		}{paidResult, allResult, nil}

		orders := make([]models.LingjiOrderUser, 0)
		session := engine.Table(tableName+"_"+strconv.Itoa(index)).Alias("o").Join("LEFT", []string{"linghit_user", "u"}, "o.account=u.account").Where("o.order_time between ? and ?", startTimestamp, endTimestamp).Where("u.account is not null")
		session.In("o.order_app", []int{2000, 438626245})
		if !allChannels {
			session.In("o.channel", channels)
		}
		if !allProducts {
			session.In("o.product_name", products)
		}
		if !allTerminalTypes {
			session.In("o.terminal_type", terminalTypes)
		}
		result.err = session.Cols("o.pay_time", "u.account", "u.gender").Find(&orders)
		if result.err == nil {
			var paidAlready = make(map[string]bool)
			var allAlready = make(map[string]bool)

			for _, order := range orders {
				_, ok := allAlready[order.LingjiOrder.Account]
				if !ok {
					allAlready[order.LingjiOrder.Account] = true
					allResult.All++
					switch order.Gender {
					case 0:
						allResult.Male++
					case 1:
						allResult.Female++
					default:
						allResult.Unkown++
					}
				}

				if order.PayTime > 0 {
					_, ok := paidAlready[order.LingjiOrder.Account]
					if !ok {
						paidAlready[order.LingjiOrder.Account] = true
						paidResult.All++
						switch order.Gender {
						case 0:
							paidResult.Male++
						case 1:
							paidResult.Female++
						default:
							paidResult.Unkown++
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
			c.errorResponse(w, 500, "获取失败，请重试")
			return
		}

		paidGenderResult.Add(tmpResult.paidGenderResult)
		allGenderResult.Add(tmpResult.allGenderResult)
	}

	//计算转化率数据
	conversion := calculateGenderConversion(paidGenderResult, allGenderResult)
	c.jsonResponse(w, map[string]interface{}{"paid": paidGenderResult, "all": allGenderResult, "conversion": conversion})
}

func (c *UserPortraitController) zxcsGender(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
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

	var paidGenderResult = new(GenderResult)
	var allGenderResult = new(GenderResult)
	var ch = make(chan *struct {
		paidGenderResult *GenderResult
		allGenderResult  *GenderResult
		err              error
	}, 8)

	//开8个协程处理
	work := func(index int) {
		paidResult := new(GenderResult)
		allResult := new(GenderResult)
		result := &struct {
			paidGenderResult *GenderResult
			allGenderResult  *GenderResult
			err              error
		}{paidResult, allResult, nil}
		orders := make([]models.ZxcsOrder, 0)
		session := engine.Table(tableName+"_"+strconv.Itoa(index)).Where("order_time between ? and ?", startTimestamp, endTimestamp)
		if !allChannels {
			session.In("channel", channels)
		}
		if !allProducts {
			session.In("product_name", products)
		}
		result.err = session.Cols("email", "gender", "pay_time").Find(&orders)
		if result.err == nil {
			paidAlready := make(map[string]bool)
			allAlready := make(map[string]bool)

			for _, order := range orders {
				_, ok := allAlready[order.Email]
				if !ok {
					allAlready[order.Email] = true
					allResult.All++
					switch order.Gender {
					case 0:
						allResult.Female++
					case 1:
						allResult.Male++
					default:
						allResult.Unkown++
					}
				}

				if order.PayTime > 0 {
					_, ok := paidAlready[order.Email]
					if !ok {
						paidAlready[order.Email] = true
						paidResult.All++
						switch order.Gender {
						case 0:
							paidResult.Female++
						case 1:
							paidResult.Male++
						default:
							paidResult.Unkown++
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
			c.errorResponse(w, 500, "获取失败，请重试")
			return
		}

		paidGenderResult.Add(tmpResult.paidGenderResult)
		allGenderResult.Add(tmpResult.allGenderResult)
	}

	//计算转化率数据
	conversion := calculateGenderConversion(paidGenderResult, allGenderResult)
	c.jsonResponse(w, map[string]interface{}{"paid": paidGenderResult, "all": allGenderResult, "conversion": conversion})
}

func (c *UserPortraitController) shopGender(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {

}

func (c *UserPortraitController) lingjiAge(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
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
	var paidAgeResult AgeResult = make(map[string]int)
	var allAgeResult AgeResult = make(map[string]int)
	var ch = make(chan *struct {
		paidAgeResult AgeResult
		allAgeResult  AgeResult
		err           error
	}, 8)

	//开8个协程处理
	work := func(index int) {
		var paidResult AgeResult = make(map[string]int)
		var allResult AgeResult = make(map[string]int)
		result := &struct {
			paidAgeResult AgeResult
			allAgeResult  AgeResult
			err           error
		}{paidResult, allResult, nil}

		orders := make([]models.LingjiOrderUser, 0)
		session := engine.Table(tableName+"_"+strconv.Itoa(index)).Alias("o").Join("LEFT", []string{"linghit_user", "u"}, "o.account=u.account").Where("o.order_time between ? and ?", startTimestamp, endTimestamp).Where("u.account is not null")
		session.In("o.order_app", []int{2000, 438626245})
		if !allChannels {
			session.In("o.channel", channels)
		}
		if !allProducts {
			session.In("o.product_name", products)
		}
		if !allTerminalTypes {
			session.In("o.terminal_type", terminalTypes)
		}
		result.err = session.Cols("o.pay_time", "u.account", "u.birthday").Find(&orders)
		if result.err == nil {
			var paidAlready = make(map[string]bool)
			var allAlready = make(map[string]bool)

			for _, order := range orders {
				_, ok := allAlready[order.LingjiOrder.Account]
				if !ok {
					allAlready[order.LingjiOrder.Account] = true
					allResult["all"]++
					if order.Birthday == 0 || int64(order.Birthday) > time.Now().Unix() {
						allResult["unkown"]++
					} else {
						age := int((time.Now().Unix() - int64(order.Birthday)) / 31536000)
						allResult[strconv.Itoa(age)]++
					}
				}

				if order.PayTime > 0 {
					_, ok := paidAlready[order.LingjiOrder.Account]
					if !ok {
						paidAlready[order.LingjiOrder.Account] = true
						paidResult["all"]++
						if order.Birthday == 0 || int64(order.Birthday) > time.Now().Unix() {
							paidResult["unkown"]++
						} else {
							age := int((time.Now().Unix() - int64(order.Birthday)) / 31536000)
							paidResult[strconv.Itoa(age)]++
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
			c.errorResponse(w, 500, "获取失败，请重试")
			return
		}

		paidAgeResult.Add(tmpResult.paidAgeResult)
		allAgeResult.Add(tmpResult.allAgeResult)
	}

	//计算转化率数据
	conversion := calculateAgeConversion(paidAgeResult, allAgeResult)
	c.jsonResponse(w, map[string]interface{}{"paid": paidAgeResult, "all": allAgeResult, "conversion": conversion})
}

func (c *UserPortraitController) zxcsAge(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
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

	var paidAgeResult AgeResult = make(map[string]int)
	var allAgeResult AgeResult = make(map[string]int)
	var ch = make(chan *struct {
		paidAgeResult AgeResult
		allAgeResult  AgeResult
		err           error
	}, 8)

	//开8个协程处理
	work := func(index int) {
		var paidResult AgeResult = make(map[string]int)
		var allResult AgeResult = make(map[string]int)
		result := &struct {
			paidAgeResult AgeResult
			allAgeResult  AgeResult
			err           error
		}{paidResult, allResult, nil}
		orders := make([]models.ZxcsOrder, 0)
		session := engine.Table(tableName+"_"+strconv.Itoa(index)).Where("order_time between ? and ?", startTimestamp, endTimestamp)
		if !allChannels {
			session.In("channel", channels)
		}
		if !allProducts {
			session.In("product_name", products)
		}
		result.err = session.Cols("email", "birthday", "pay_time").Find(&orders)
		if result.err == nil {
			paidAlready := make(map[string]bool)
			allAlready := make(map[string]bool)

			for _, order := range orders {
				_, ok := allAlready[order.Email]
				if !ok {
					allAlready[order.Email] = true
					allResult["all"]++
					if order.Birthday == 0 || int64(order.Birthday) > time.Now().Unix() {
						allResult["unkown"]++
					} else {
						age := int((time.Now().Unix() - int64(order.Birthday)) / 31536000)
						allResult[strconv.Itoa(age)]++
					}
				}

				if order.PayTime > 0 {
					_, ok := paidAlready[order.Email]
					if !ok {
						paidAlready[order.Email] = true
						paidResult["all"]++
						if order.Birthday == 0 || int64(order.Birthday) > time.Now().Unix() {
							paidResult["unkown"]++
						} else {
							age := int((time.Now().Unix() - int64(order.Birthday)) / 31536000)
							paidResult[strconv.Itoa(age)]++
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
			c.errorResponse(w, 500, "获取失败，请重试")
			return
		}

		paidAgeResult.Add(tmpResult.paidAgeResult)
		allAgeResult.Add(tmpResult.allAgeResult)
	}

	//计算消费转化率
	conversion := calculateAgeConversion(paidAgeResult, allAgeResult)
	c.jsonResponse(w, map[string]interface{}{"paid": paidAgeResult, "all": allAgeResult, "conversion": conversion})
}

func (c *UserPortraitController) shopAge(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {

}

func (c *UserPortraitController) lingjiHour(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
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
	var paidHourResult HourResult = make(map[string]int)
	var allHourResult HourResult = make(map[string]int)
	var ch = make(chan *struct {
		paidHourResult HourResult
		allHourResult  HourResult
		err            error
	}, 8)

	//开8个协程处理
	work := func(index int) {
		var paidResult HourResult = make(map[string]int)
		var allResult HourResult = make(map[string]int)
		result := &struct {
			paidHourResult HourResult
			allHourResult  HourResult
			err            error
		}{paidResult, allResult, nil}

		orders := make([]models.LingjiOrder, 0)
		session := engine.Table(tableName+"_"+strconv.Itoa(index)).Where("order_time between ? and ?", startTimestamp, endTimestamp)
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
		result.err = session.Cols("pay_time", "order_time", "account", "device_id").Find(&orders)
		if result.err == nil {
			var paidAlready = make(map[string][]bool)
			var allAlready = make(map[string][]bool)

			for _, order := range orders {
				//获取订单时间点
				time := time.Unix(int64(order.OrderTime), 0).Local()
				hour := time.Hour()

				//获取用户唯一表示
				var identifier string
				if order.Account != "" {
					identifier = "account:" + order.Account
				} else if order.DeviceId != "" {
					identifier = "device_id:" + order.DeviceId
				} else {
					continue
				}

				if _, ok := allAlready[identifier]; !ok {
					allAlready[identifier] = make([]bool, 24)
				}
				if !allAlready[identifier][hour] {
					allAlready[identifier][hour] = true
					allResult["all"]++
					allResult[strconv.Itoa(hour)]++
				}

				if order.PayTime > 0 {
					if _, ok := paidAlready[identifier]; !ok {
						paidAlready[identifier] = make([]bool, 24)
					}
					if !paidAlready[identifier][hour] {
						paidAlready[identifier][hour] = true
						paidResult["all"]++
						paidResult[strconv.Itoa(hour)]++
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
			c.errorResponse(w, 500, "获取失败，请重试")
			return
		}

		paidHourResult.Add(tmpResult.paidHourResult)
		allHourResult.Add(tmpResult.allHourResult)
	}

	//计算消费转化率
	conversion := calculateHourConversion(paidHourResult, allHourResult)
	c.jsonResponse(w, map[string]interface{}{"paid": paidHourResult, "all": allHourResult, "conversion": conversion})
}

func (c *UserPortraitController) zxcsHour(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {
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

	var paidHourResult HourResult = make(map[string]int)
	var allHourResult HourResult = make(map[string]int)
	var ch = make(chan *struct {
		paidHourResult HourResult
		allHourResult  HourResult
		err            error
	}, 8)

	//开8个协程处理
	work := func(index int) {
		var paidResult HourResult = make(map[string]int)
		var allResult HourResult = make(map[string]int)
		result := &struct {
			paidHourResult HourResult
			allHourResult  HourResult
			err            error
		}{paidResult, allResult, nil}
		orders := make([]models.ZxcsOrder, 0)
		session := engine.Table(tableName+"_"+strconv.Itoa(index)).Where("order_time between ? and ?", startTimestamp, endTimestamp)
		if !allChannels {
			session.In("channel", channels)
		}
		if !allProducts {
			session.In("product_name", products)
		}
		result.err = session.Cols("email", "pay_time", "order_time").Find(&orders)
		if result.err == nil {
			paidAlready := make(map[string][]bool)
			allAlready := make(map[string][]bool)

			for _, order := range orders {
				//获取订单时间点
				time := time.Unix(int64(order.OrderTime), 0).Local()
				hour := time.Hour()

				_, ok := allAlready[order.Email]
				if !ok {
					allAlready[order.Email] = make([]bool, 24)
				}
				if !allAlready[order.Email][hour] {
					allAlready[order.Email][hour] = true

					allResult["all"]++
					allResult[strconv.Itoa(hour)]++
				}

				if order.PayTime > 0 {
					_, ok := paidAlready[order.Email]
					if !ok {
						paidAlready[order.Email] = make([]bool, 24)
					}
					if !paidAlready[order.Email][hour] {
						paidAlready[order.Email][hour] = true

						paidResult["all"]++
						paidResult[strconv.Itoa(hour)]++
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
			c.errorResponse(w, 500, "获取失败，请重试")
			return
		}

		paidHourResult.Add(tmpResult.paidHourResult)
		allHourResult.Add(tmpResult.allHourResult)
	}

	//计算消费转化率
	conversion := calculateHourConversion(paidHourResult, allHourResult)
	c.jsonResponse(w, map[string]interface{}{"paid": paidHourResult, "all": allHourResult, "conversion": conversion})
}

func (c *UserPortraitController) shopHour(w http.ResponseWriter, r *http.Request, data map[string]interface{}) {

}
