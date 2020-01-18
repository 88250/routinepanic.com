// 协慌网 - 专注编程问答汉化 https://routinepanic.com
// Copyright (C) 2018-present, b3log.org
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package service

import (
	"context"
	"strings"

	"cloud.google.com/go/translate"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/language"
)

// Translation service.
var Translation = &translationService{}

type translationService struct {
}

func (srv *translationService) Translate(text string, format string) string {
	//dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:1081", nil, proxy.Direct)
	//if err != nil {
	//	logger.Fatal("can't connect to the proxy: " + err.Error())
	//
	//	return ""
	//}
	//
	//httpTransport := &http.Transport{Dial: dialer.Dial}
	//httpTransport := &http.Transport{Dial: dialer.Dial}
	//http.DefaultClient.Transport = httpTransport

	ctx := context.Background()
	client, err := translate.NewClient(ctx)
	if err != nil {
		logger.Errorf("create translate client failed: " + err.Error())

		return ""
	}

	ret := ""

	translations, err := client.Translate(ctx, []string{text}, language.Chinese,
		&translate.Options{Source: language.English, Format: translate.Format(format), Model: "nmt"})
	if nil == err {
		ret = translations[0].Text
	}

	if "" == ret {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(text))
		if nil != err {
			logger.Errorf("parse text to HTML doc failed: " + err.Error())

			return ""
		}

		fragment := ""
		pCount := 0
		doc.Find("body").Children().Each(func(i int, s *goquery.Selection) {
			nodeName := goquery.NodeName(s)
			html, _ := s.Html()
			if "pre" == nodeName || "code" == nodeName {
				ret += translateFragment(client, ctx, fragment)
				ret += "<" + nodeName + ">" + html + "</" + nodeName + ">"
				fragment = ""
				pCount = 0

				return
			}

			if "" == html {
				fragment += "<" + nodeName + ">"
			} else {
				fragment += "<" + nodeName + ">" + html + "</" + nodeName + ">"
			}

			if "p" == nodeName {
				pCount++
			}

			if 3 < pCount {
				ret += translateFragment(client, ctx, fragment)
				fragment = ""
				pCount = 0
			}
		})

		if "" != fragment {
			ret += translateFragment(client, ctx, fragment)
		}
	}

	return ret
}

func translateFragment(client *translate.Client, ctx context.Context, fragment string) string {
	translations, err := client.Translate(ctx, []string{fragment}, language.Chinese,
		&translate.Options{Source: language.English, Format: translate.HTML, Model: "nmt"})
	if nil != err {
		logger.Errorf("translate failed: " + err.Error())

		return fragment
	}

	translated := translations[0].Text
	if "" != translated {
		return translated
	}

	return fragment
}
