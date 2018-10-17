package util

import (
	"time"

	"github.com/parnurzeal/gorequest"
)

func PushBaidu(urls string) {
	if "" == Conf.BaiduToken {
		return
	}

	baiduURL := "http://data.zz.baidu.com/urls?site=" + Conf.Server + "&token=" + Conf.BaiduToken
	logger.Info(baiduURL)
	_, data, errors := gorequest.New().Post(baiduURL).
		AppendHeader("User-Agent", "curl/7.12.1").
		AppendHeader("Host", "data.zz.baidu.com").
		AppendHeader("Content-Type", "text/plain").Timeout(10 * time.Second).Send(urls).EndBytes()
	if nil != errors {
		logger.Errorf("push to baidu failed [%s]", errors)

		return
	}

	logger.Info(string(data))
}
