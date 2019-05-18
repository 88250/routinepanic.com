// 协慌网 https://routinepanic.com
// Copyright (C) 2018-present, b3log.org

package util

import (
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"
)

func PushBaidu(urls string) {
	if "" == Conf.BaiduToken {
		return
	}

	if !strings.HasSuffix(urls, "\n") {
		urls += "\n"
	}

	site := Conf.Server
	site = site[strings.Index(site, "://")+len("://"):]
	baiduURL := "http://data.zz.baidu.com/urls?site=" + site + "&token=" + Conf.BaiduToken
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
