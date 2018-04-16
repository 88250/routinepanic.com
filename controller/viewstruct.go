// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package controller

import "html/template"

type question struct {
	ID    uint64
	Path  string
	Title string
	Tags  []*tag
	Content template.HTML
}

type tag struct {
	Title string
}
