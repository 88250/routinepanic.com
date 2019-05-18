// 协慌网 https://routinepanic.com
// Copyright (C) 2018-present, b3log.org

package controller

import (
	"errors"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/b3log/routinepanic.com/log"
	"github.com/b3log/routinepanic.com/util"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
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

	store := cookie.NewStore([]byte(util.Conf.SessionSecret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   util.Conf.SessionMaxAge,
		Secure:   strings.HasPrefix(util.Conf.Server, "https"),
		HttpOnly: true,
	})
	ret.Use(sessions.Sessions("rp", store), fillCommon)

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
	ret.StaticFile("/robots.txt", "view/robots.txt")

	ret.GET("", showIndexAction)
	ret.GET("/questions/:path", showQuestionAction)
	ret.GET("/questions/:path/answers/:answerID", showQuestionAnswerAction)
	ret.GET("/tags/*tag", showTagAction)
	ret.GET("/baidu", submitURL)
	ret.GET("/import/so", importSO)
	ret.GET("/replenish/answers", replenishAnswers)
	ret.GET("/translate", translate)

	ret.GET("/contri/:dataType/:id", showContriAction)
	ret.POST("/contri/:dataType/:id", contriAction)
	ret.POST("/html", tuneHTMLAction)

	ret.GET("/words/:name", getWordAction)

	ret.GET("/reviews/waiting", showWaitingReviewsAction)
	ret.GET("/reviews/passed", showPassedReviewsAction)
	ret.GET("/reviews/rejected", showRejectedReviewsAction)
	ret.GET("/reviews/details/:id", showReviewAction)
	ret.POST("/reviews/review/:id", ReviewAction)

	ret.GET("/login", showLoginAction)
	ret.GET("/logout", logoutAction)
	ret.GET("/oauth/github/redirect", redirectGitHubAction)
	ret.GET("/oauth/github/callback", githubCallbackHandler)

	return ret
}

// DataModel represents data model.
type DataModel map[string]interface{}

func (dataModel *DataModel) Put(key string, value interface{}) {
	(*dataModel)[key] = value
}

func (dataModel *DataModel) GetStr(key string) string {
	return (*dataModel)[key].(string)
}

func fillCommon(c *gin.Context) {
	dataModel := &DataModel{}
	c.Set("dataModel", dataModel)

	dataModel.Put("Conf", util.Conf)
	dataModel.Put("Title", "协慌网")
	dataModel.Put("Slogan", util.Slogan)
	dataModel.Put("MetaKeywords", util.MetaKeywords)
	dataModel.Put("MetaDescription", util.Slogan)
	session := util.GetSession(c)
	dataModel.Put("User", session)
}

func getDataModel(c *gin.Context) *DataModel {
	dataModelVal, _ := c.Get("dataModel")

	return dataModelVal.(*DataModel)
}

func notFound(c *gin.Context) {
	dataModel := getDataModel(c)
	c.HTML(http.StatusNotFound, "404.html", dataModel)
}
