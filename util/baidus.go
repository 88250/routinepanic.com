package util

import (
	"net/url"
	"time"

	"github.com/parnurzeal/gorequest"
)

func PushBaidu(urls string) {
	if "" == Conf.BaiduToken {
		return
	}

	baiduURL := "http://data.zz.baidu.com/urls?site=" + url.QueryEscape(Conf.Server) +
		"&token=" + url.QueryEscape(Conf.BaiduToken)
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
