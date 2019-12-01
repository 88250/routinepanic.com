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
	"crypto/tls"
	"github.com/88250/gulu"
	"net/http"
	"strings"
	"time"

	"github.com/88250/routinepanic.com/model"
	"github.com/88250/routinepanic.com/service"
	"github.com/88250/routinepanic.com/util"
	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
)

var states = map[string]string{}

func redirectGitHubAction(c *gin.Context) {
	requestResult := gulu.Ret.NewResult()
	_, _, errs := gorequest.New().TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		Get("https://hacpai.com/oauth/rp/client").
		Set("user-agent", util.UserAgent).Timeout(10 * time.Second).EndStruct(requestResult)
	if nil != errs {
		logger.Errorf("Get oauth client id failed: %+v", errs)
		c.Status(http.StatusInternalServerError)

		return
	}
	if 0 != requestResult.Code {
		logger.Errorf("get oauth client id failed [code=%d, msg=%s]", requestResult.Code, requestResult.Msg)
		c.Status(http.StatusNotFound)

		return
	}
	data := requestResult.Data.(map[string]interface{})
	clientId := data["clientId"].(string)
	loginAuthURL := data["loginAuthURL"].(string)

	state := c.Query("state")
	referer := util.Conf.Server + "__" + state
	state = gulu.Rand.String(16) + referer
	states[state] = state
	path := loginAuthURL + "?client_id=" + clientId + "&state=" + state
	c.Redirect(http.StatusSeeOther, path)
}

func githubCallbackHandler(c *gin.Context) {
	state := c.Query("state")
	if _, exist := states[state]; !exist {
		c.Status(http.StatusBadRequest)

		return
	}
	delete(states, state)

	referer := state[16:]
	if strings.Contains(referer, "__0") || strings.Contains(referer, "__1") {
		referer = referer[:len(referer)-len("__0")]
	}
	accessToken := c.Query("ak")
	githubUser := GitHubUserInfo(accessToken)
	if nil == githubUser {
		logger.Warnf("Can not get user info with token [" + accessToken + "]")
		c.Status(http.StatusUnauthorized)

		return
	}

	githubId := githubUser["userId"].(string)
	userName := githubUser["userName"].(string)
	avatar := githubUser["userAvatarURL"].(string)
	user := &model.User{Name: userName, Avatar: avatar, GithubId: githubId}
	if err := service.User.AddOrUpdate(user); nil != err {
		logger.Errorf("add a new user failed: " + err.Error())
		c.Status(http.StatusInternalServerError)

		return
	}

	role := util.RoleNormal
	if "88250" == userName {
		role = util.RoleReviewer
	}

	session := util.SessionData{UID: user.ID, UName: userName, UAvatar: avatar, URole: role}
	session.Save(c)

	c.Redirect(http.StatusSeeOther, "/")
}

// GitHubUserInfo returns GitHub user info specified by the given access token.
func GitHubUserInfo(accessToken string) (ret map[string]interface{}) {
	result := map[string]interface{}{}
	response, data, errors := gorequest.New().TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		Get("https://hacpai.com/github/user?ak="+accessToken).Timeout(7*time.Second).
		Set("User-Agent", util.UserAgent).EndStruct(&result)
	if nil != errors || http.StatusOK != response.StatusCode {
		logger.Errorf("Get github user info failed: %+v, %s", errors, data)

		return nil
	}

	if 0 != result["sc"].(float64) {
		return nil
	}

	return result["data"].(map[string]interface{})
}

func showLoginAction(c *gin.Context) {
	dataModel := getDataModel(c)
	c.HTML(http.StatusOK, "login.html", dataModel)
}

func logoutAction(c *gin.Context) {
	util.Invalidate(c)
	c.Redirect(http.StatusSeeOther, "/")
}
