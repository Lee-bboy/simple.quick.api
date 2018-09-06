package common

import (
	"time"

	"datacenter.analysis.api/config"
)

var Location *time.Location

func init() {
	var err error
	Location, err = time.LoadLocation(GetConf("timezone", "Asia/Shanghai"))
	if err != nil {
		panic(err)
	}
}

func GetConf(key string, value string) (data string) {
	val, err := config.Config.Get(key)

	if err != nil {
		val = value
	}

	data = val

	return data
}

//获取用户的accessSecret
func GetAccessSecret(accessKey string) (string, error) {
	accessSecret, err := config.Config.Get("auth." + accessKey)
	if err != nil {
		return "", err
	}

	return accessSecret, nil
}
