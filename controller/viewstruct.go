// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package controller

type question struct {
	ID    uint64
	Path  string
	Title string
	Tags  []*tag
}

type tag struct {
	Title string
}
