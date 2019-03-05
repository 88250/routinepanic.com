// 协慌网 https://routinepanic.com
// Copyright (C) 2018-2019, b3log.org

package controller

import (
	"github.com/b3log/routinepanic.com/util"
	"github.com/gin-gonic/gin"
)

func fillUser(c *gin.Context) {
	session := util.GetSession(c)
	dataModel := getDataModel(c)
	dataModel.Put("User", session)
	if "" != session.UName {
		c.Next()

		return
	}

	uaStr := c.Request.UserAgent()
	logger.Tracef("Bot User-Agent [%s]", uaStr)
	c.Next()
}
