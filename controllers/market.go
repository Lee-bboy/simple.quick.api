package controllers

import (
	"net/http"
	"strconv"

	"time"

	"datacenter.analysis.api/common"
	"datacenter.analysis.api/helper"
)

var MDController = new(MarketDataController)

type MarketDataController struct {
	Controller
}

func (c *MarketDataController) Associate(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	data, err := c.parseParams(r)

	if err != nil {
		c.errorResponse(w, 400, "请求参数错误")
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
	} else {
		c.errorResponse(w, 400, "产品不可以为空")
		return
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

	var res interface{}

	switch businessId {
	case 1:
		res, err = helper.LingjiAssociate(data)
		if err != nil {
			c.errorResponse(w, 500, "获取失败，请重试")
			return
		}
	case 2:
		res = helper.ZxcsAssociate(data)
	case 3:
		//helper.ShopAssociate(data)
	default:
		c.errorResponse(w, 400, "业务线找不到")
		return
	}

	c.jsonResponse(w, res)
}

func (c *MarketDataController) Frequency(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	data, err := c.parseParams(r)
	if err != nil {
		c.errorResponse(w, 400, "请求参数错误")
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

	var res interface{}

	switch businessId {
	case 1:
		res, err = helper.LingjiFrequency(data)
		if err != nil {
			c.errorResponse(w, 500, "获取失败，请重试")
			return
		}
	case 2:
		res = helper.ZxcsFrequency(data)
	case 3:
		helper.ShopFrequency(data)
	default:
		c.errorResponse(w, 400, "业务线找不到")
		return
	}

	c.jsonResponse(w, res)
}

func (c *MarketDataController) TimeInterval(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	data, err := c.parseParams(r)
	if err != nil {
		c.errorResponse(w, 400, "请求参数错误")
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

	var res interface{}

	switch businessId {
	case 1:
		res, err = helper.LingjiTimeInterval(data)
		if err != nil {
			c.errorResponse(w, 500, "获取失败，请重试")
			return
		}
	case 2:
		res = helper.ZxcsTimeInterval(data)
	case 3:
		//helper.ShopTimeInterval(data)
	default:
		c.errorResponse(w, 400, "业务线找不到")
		return
	}

	c.jsonResponse(w, res)
}

func (c *MarketDataController) FirstTrade(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	data, err := c.parseParams(r)
	if err != nil {
		c.errorResponse(w, 400, "请求参数错误")
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

	var res interface{}

	switch businessId {
	case 1:
		res, err = helper.LingjiFirstTrade(data)
		if err != nil {
			c.errorResponse(w, 500, "获取失败，请重试")
			return
		}
	case 2:
		//res = helper.ZxcsTimeInterval(data)
	case 3:
		//helper.ShopTimeInterval(data)
	default:
		c.errorResponse(w, 400, "业务线找不到")
		return
	}

	c.jsonResponse(w, res)
}
