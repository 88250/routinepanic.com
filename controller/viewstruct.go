// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package controller

type question struct {
	Title string
	Tags []*tag
}

type tag struct {
	Title string
}
