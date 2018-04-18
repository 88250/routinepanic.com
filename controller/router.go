// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package controller

import (
	"errors"
	"html/template"
	"os"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/b3log/routinepanic.com/log"
	"github.com/b3log/routinepanic.com/util"
	"path/filepath"
	"net/http"
)

// Logger
var logger = log.NewLogger(os.Stdout)

// MapRoutes returns a gin engine and binds controllers with request URLs.
func MapRoutes() *gin.Engine {
	ret := gin.New()
	ret.SetFuncMap(template.FuncMap{
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("len(values) is " + strconv.Itoa(len(values)%2))
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, errors.New("")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
		"minus": func(a, b int) int {
			return a - b
		},
	})

	ret.Use(gin.Recovery())

	store := sessions.NewCookieStore([]byte(util.Conf.SessionSecret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   util.Conf.SessionMaxAge,
		Secure:   strings.HasPrefix(util.Conf.Server, "https"),
		HttpOnly: true,
	})
	ret.Use(sessions.Sessions("rp", store))

	templates, err := filepath.Glob("view/template/*.html")
	if nil != err {
		logger.Fatal("load templates failed: " + err.Error())
	}
	subTemplates, _ := filepath.Glob("view/template/*/*.html")
	templates = append(templates, subTemplates...)
	ret.LoadHTMLFiles(templates...)

	ret.Static("/css", "view/css")
	ret.Static("/js", "view/js")
	ret.Static("/images", "view/images")

	ret.Use(fillCommon)
	ret.GET("", showIndexAction)
	ret.GET("/questions/*path", showQuestionAction)

	return ret
}

// DataModel represents data model.
type DataModel map[string]interface{}

func fillCommon(c *gin.Context) {
	dataModel := &DataModel{}
	c.Set("dataModel", dataModel)

	(*dataModel)["Conf"] = util.Conf
	(*dataModel)["Title"] = "协慌网"
	(*dataModel)["Slogan"] = util.Slogan
}

func getDataModel(c *gin.Context) DataModel {
	dataModelVal, _ := c.Get("dataModel")

	return *(dataModelVal.(*DataModel))
}

func notFound(c *gin.Context) {
	dataModel := getDataModel(c)
	c.HTML(http.StatusNotFound, "404.html", dataModel)
}