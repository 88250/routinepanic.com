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

	"github.com/b3log/routinepanic.com/service"
	"github.com/b3log/routinepanic.com/util"
	"github.com/gin-gonic/gin"
)

func submitURL(c *gin.Context) {
	for i := 1; i < 256; i++ {
		questions, pagination := service.QnA.GetQuestions(i)
		if i > pagination.LastPageNum {
			logger.Infof("submit completed [p=%d]", pagination.LastPageNum)

			break
		}

		urls := ""
		for _, question := range questions {
			urls += util.Conf.Server + "/questions/" + question.Path + "\n"
		}

		util.PushBaidu(urls)
	}

	c.Redirect(http.StatusTemporaryRedirect, util.Conf.Server)
}
