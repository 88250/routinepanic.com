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

package util

import (
	"encoding/json"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	RoleNormal = iota
	RoleReviewer
)

// SessionData represents the session.
type SessionData struct {
	UID     uint64 // user id
	UName   string // username
	UAvatar string // user avatar URL
	URole   int    // 0: normal user, 1: reviewer
}

// Save saves the current session of the specified context.
func (sd *SessionData) Save(c *gin.Context) error {
	session := sessions.Default(c)
	sessionDataBytes, err := json.Marshal(sd)
	if nil != err {
		return err
	}
	session.Set("data", string(sessionDataBytes))

	return session.Save()
}

func IsLoggedIn(c *gin.Context) bool {
	session := GetSession(c)

	return "" != session.UName
}

// Invalidate invalidates the current session.
func Invalidate(c *gin.Context) {
	session := sessions.Default(c)
	session.Options(sessions.Options{Path: "/", MaxAge: -1})
	session.Clear()
	session.Save()
}

// GetSession returns session of the specified context.
func GetSession(c *gin.Context) *SessionData {
	ret := &SessionData{}

	session := sessions.Default(c)
	sessionDataStr := session.Get("data")
	if nil == sessionDataStr {
		return ret
	}

	err := json.Unmarshal([]byte(sessionDataStr.(string)), ret)
	if nil != err {
		return ret
	}

	c.Set("session", ret)

	return ret
}
