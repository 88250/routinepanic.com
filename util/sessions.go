// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package util

import (
	"encoding/json"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// SessionData represents the session.
type SessionData struct {
	UName   string // username
	UAvatar string // user avatar URL
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
