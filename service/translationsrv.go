// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package service

import (
	"context"
	"net/http"

	"cloud.google.com/go/translate"
	"golang.org/x/net/proxy"
	"golang.org/x/text/language"
)

// Translation service.
var Translation = &translationService{}

type translationService struct {
}

func (srv *translationService) Translate(text string, format string) string {
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:1081", nil, proxy.Direct)
	if err != nil {
		logger.Fatal("can't connect to the proxy: " + err.Error())
	}

	httpTransport := &http.Transport{Dial: dialer.Dial}
	http.DefaultClient.Transport = httpTransport

	ctx := context.Background()
	client, err := translate.NewClient(ctx)
	if err != nil {
		logger.Errorf("create translate client failed: " + err.Error())

		return ""
	}

	translations, err := client.Translate(ctx, []string{text}, language.Chinese,
		&translate.Options{Source: language.English, Format: translate.Format(format), Model: "nmt"})
	if nil != err {
		logger.Errorf("translate failed: " + err.Error())

		return ""
	}

	return translations[0].Text
}
