// 协慌网 https://routinepanic.com
// Copyright (C) 2018-present, b3log.org

// 协慌网 https://routinepanic.com
// Copyright (C) 2018-2019, b3log.org

package util

// Result represents HTTP response body.
type Result struct {
	Code int         `json:"code"` // return code, 0 for succ
	Msg  string      `json:"msg"`  // message
	Data interface{} `json:"data"` // data object
}

// NewResult creates a result with Code=0, Msg="", Data=nil.
func NewResult() *Result {
	return &Result{
		Code: 0,
		Msg:  "",
		Data: nil,
	}
}
