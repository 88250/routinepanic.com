// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package controller

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/b3log/routinepanic.com/util"
	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	"github.com/parnurzeal/gorequest"
)

const nilB3id = "H9oxzSym"

func fillUser(c *gin.Context) {
	dataModel := &DataModel{}
	c.Set("dataModel", dataModel)
	session := util.GetSession(c)
	(*dataModel)["User"] = session
	if 0 != session.UID {
		c.Next()

		return
	}

	uaStr := c.Request.UserAgent()
	if isBot(uaStr) {
		logger.Tracef("Bot User-Agent [%s]", uaStr)
		c.Next()

		return
	}

	b3id := c.Request.URL.Query().Get("b3id")
	switch b3id {
	case nilB3id:
		c.Next()

		return
	case "":
		redirectURL := util.Conf.Server + c.Request.URL.Path
		redirectURL = url.QueryEscape(redirectURL)
		c.Redirect(http.StatusSeeOther, "https://hacpai.com/apis/b3-identity?goto="+redirectURL)
		c.Abort()

		return
	default:
		result := util.NewResult()
		_, _, errs := gorequest.New().Get("https://hacpai.com/apis/check-b3-identity?b3id="+b3id).
			Set("user-agent", util.UserAgent).Timeout(5*time.Second).
			Retry(3, 2*time.Second, http.StatusInternalServerError).EndStruct(result)
		if nil != errs {
			logger.Errorf("check b3 identity failed: %s", errs)
			c.Next()

			return
		}

		if 0 != result.Code {
			c.Next()

			return
		}

		data := result.Data.(map[string]interface{})
		username := data["userName"].(string)
		userAvatar := data["userAvatarURL"].(string)

		session = &util.SessionData{
			UName:   username,
			UAvatar: userAvatar,
		}

		if err := session.Save(c); nil != err {
			result.Code = -1
			result.Msg = "saves session failed: " + err.Error()
		}

		(*dataModel)["User"] = session
		c.Next()
	}
}

func isBot(uaStr string) bool {
	var ua = user_agent.New(uaStr)

	return ua.Bot() || strings.HasPrefix(uaStr, "Sym")
}