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

	"github.com/88250/gulu"
	"github.com/88250/routinepanic.com/model"
	"github.com/88250/routinepanic.com/service"
	"github.com/88250/routinepanic.com/util"
	"github.com/gin-gonic/gin"
)

var states = map[string]string{}

func redirectLoginAction(c *gin.Context) {
	loginAuthURL := "https://ld246.com/login?goto=" + util.Conf.Server + "/login/callback"

	state := gulu.Rand.String(16)
	states[state] = state
	path := loginAuthURL + "&state=" + state + "&v=" + util.Version
	c.Redirect(http.StatusSeeOther, path)
}

func loginCallbackHandler(c *gin.Context) {
	state := c.Query("state")
	if _, exist := states[state]; !exist {
		c.Status(http.StatusBadRequest)
		return
	}
	delete(states, state)

	accessToken := c.Query("access_token")
	userInfo := util.HacPaiUserInfo(accessToken)

	userId := userInfo["userId"].(string)
	userName := userInfo["userName"].(string)
	avatar := userInfo["avatar"].(string)
	user := &model.User{Name: userName, Avatar: avatar, GithubId: userId}
	if err := service.User.AddOrUpdate(user); nil != err {
		logger.Errorf("add or update a user failed: " + err.Error())
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

func showLoginAction(c *gin.Context) {
	dataModel := getDataModel(c)
	dataModel.Put("AdSense", false)
	c.HTML(http.StatusOK, "login.html", dataModel)
}

func logoutAction(c *gin.Context) {
	util.Invalidate(c)
	c.Redirect(http.StatusSeeOther, "/")
}
