// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package service

import (
	"net/http"
	"testing"

	"golang.org/x/net/proxy"
)

func TestTranslate(t *testing.T) {
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:1081", nil, proxy.Direct)
	if err != nil {
		t.Fatal("can't connect to the proxy: " + err.Error())
	}

	httpTransport := &http.Transport{Dial: dialer.Dial}
	http.DefaultClient.Transport = httpTransport

	text := Translation.Translate("Why is it faster to process a sorted array than an unsorted array?")
	t.Log(text)
}
