// 协慌网 https://routinepanic.com
// Copyright (C) 2018-2019, b3log.org

// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

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
