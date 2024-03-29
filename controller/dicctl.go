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

package controller

import (
	"net/http"
	"strings"
	"time"

	"github.com/88250/gulu"
	"github.com/88250/routinepanic.com/model"
	"github.com/88250/routinepanic.com/service"
	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
)

func getWordAction(c *gin.Context) {
	result := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, result)

	wordName := strings.ToLower(c.Param("name"))
	word := service.Dic.GetWord(wordName)
	if nil == word {
		wordResult := map[string]interface{}{}
		response, bytes, errors := gorequest.New().Get("http://dict-co.iciba.com/api/dictionary.php").
			Param("w", wordName).
			Param("type", "json").
			Param("key", "AD4631FD78BDD04FE30000C76556C133").
			Timeout(7*time.Second).Retry(3, time.Second).EndStruct(&wordResult)
		if nil != errors {
			logger.Errorf("query iciba failed: %v", errors)
			result.Code = -1

			return
		}
		if http.StatusOK != response.StatusCode {
			logger.Errorf("query iciba failed: %d, %s", response.StatusCode, bytes)
			result.Code = -1

			return
		}

		word = parseWordResult(wordResult)
		if nil != word {
			if err := service.Dic.AddWord(word); nil != err {
				logger.Errorf("add word failed: %s", err.Error())
			}
		}
	}

	if nil == word {
		result.Code = -1

		return
	}

	result.Data = word
}

func parseWordResult(wordResult map[string]interface{}) *model.Word {
	val, ok := wordResult["symbols"]
	if !ok {
		return nil
	}

	symbols := val.([]interface{})
	if 1 > len(symbols) {
		return nil
	}

	if 1 < len(symbols) {
		logger.Infof("found a symbols, word [%+v]", wordResult)
	}

	symbol := symbols[0].(map[string]interface{})
	val, ok = symbol["parts"]
	if !ok {
		return nil
	}

	parts := val.([]interface{})
	if 1 > len(parts) {
		return nil
	}

	meansBuilder := ""
	for _, pVal := range parts {
		part := pVal.(map[string]interface{})
		val, ok = part["part"]
		if !ok {
			return nil
		}

		meansBuilder += val.(string) + " "
		means := part["means"].([]interface{})
		for i, mVal := range means {
			mean := mVal.(string)
			meansBuilder += mean
			if i < len(means)-1 {
				meansBuilder += "；"
			}
		}
		meansBuilder += "\n"
	}

	return &model.Word{
		Name:    wordResult["word_name"].(string),
		PhAm:    symbol["ph_am"].(string),
		PhAmMp3: symbol["ph_am_mp3"].(string),
		PhEn:    symbol["ph_en"].(string),
		PhEnMp3: symbol["ph_en_mp3"].(string),
		Means:   meansBuilder,
	}
}
