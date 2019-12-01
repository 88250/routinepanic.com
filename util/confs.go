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

// Package util defines variety of utilities.
package util

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/88250/gulu"
	"github.com/88250/routinepanic.com/model"
	"github.com/jinzhu/gorm"
)

// Slogan
const Slogan = "专注编程问答汉化"

// Meta keywords
const MetaKeywords = "程序员,编程,代码,问答"

// Logger
var logger = gulu.Log.NewLogger(os.Stdout)

// Version of RP.
const Version = "1.0.0"

// User-Agent of RP.
const UserAgent = "Mozilla/5.0 (compatible; RP/" + Version + "; +https://routinepanic.com)"

// Conf of RP.
var Conf *Configuration

// Models represents all models..
var Models = []interface{}{
	&model.Question{}, &model.Answer{}, &model.Tag{}, &model.Correlation{}, &model.Revision{}, &model.Word{}, &model.User{},
	&model.Review{},
}

// Pagination parameters.
const (
	WindowSize = 20
	PageSize   = 50
)

// Table prefix.
const tablePrefix = "rp_"

// Configuration (rp.json).
type Configuration struct {
	Server                string // server scheme, host and port
	StaticServer          string // static resources server scheme, host and port
	StaticResourceVersion string // version of static resources
	LogLevel              string // logging level: trace/debug/info/warn/error/fatal
	SessionSecret         string // HTTP session secret
	SessionMaxAge         int    // HTTP session max age (in seciond)
	RuntimeMode           string // runtime mode (dev/prod)
	MySQL                 string // MySQL connection URL
	Port                  string // listen port
	BaiduToken            string // baidu search push token
}

// LoadConf loads the configurations. Command-line arguments will override configuration file.
func LoadConf() {
	version := flag.Bool("version", false, "prints current version")
	confPath := flag.String("conf", "rp.json", "path of rp.json")
	confServer := flag.String("server", "", "this will override Conf.Server if specified")
	confStaticServer := flag.String("static_server", "", "this will override Conf.StaticServer if specified")
	confStaticResourceVer := flag.String("static_resource_ver", "", "this will override Conf.StaticResourceVersion if specified")
	confLogLevel := flag.String("log_level", "", "this will override Conf.LogLevel if specified")
	confRuntimeMode := flag.String("runtime_mode", "", "this will override Conf.RuntimeMode if specified")
	confMySQL := flag.String("mysql", "", "this will override Conf.MySQL if specified")
	confPort := flag.String("port", "", "this will override Conf.Port if specified")
	confBaiduToken := flag.String("baidu_token", "", "this will override Conf.BaiduToken if specified")

	flag.Parse()

	if *version {
		fmt.Println(Version)

		os.Exit(0)
	}

	bytes, err := ioutil.ReadFile(*confPath)
	if nil != err {
		logger.Fatal("loads configuration file [" + *confPath + "] failed: " + err.Error())
	}

	Conf = &Configuration{}
	if err = json.Unmarshal(bytes, Conf); nil != err {
		logger.Fatal("parses [rp.json] failed: ", err)
	}

	gulu.Log.SetLevel(Conf.LogLevel)
	if "" != *confLogLevel {
		Conf.LogLevel = *confLogLevel
		gulu.Log.SetLevel(*confLogLevel)
	}

	home, err := gulu.OS.Home()
	if nil != err {
		logger.Fatal("can't find user home directory: " + err.Error())
	}
	logger.Debugf("${home} [%s]", home)

	if "" != *confRuntimeMode {
		Conf.RuntimeMode = *confRuntimeMode
	}

	if "" != *confServer {
		Conf.Server = *confServer
	}

	if "" != *confStaticServer {
		Conf.StaticServer = *confStaticServer
	}
	if "" == Conf.StaticServer {
		Conf.StaticServer = Conf.Server
	}

	time := strconv.FormatInt(time.Now().UnixNano(), 10)
	logger.Debugf("${time} [%s]", time)
	Conf.StaticResourceVersion = strings.Replace(Conf.StaticResourceVersion, "${time}", time, 1)
	if "" != *confStaticResourceVer {
		Conf.StaticResourceVersion = *confStaticResourceVer
	}

	if "" != *confMySQL {
		Conf.MySQL = *confMySQL
	}

	if "" != *confPort {
		Conf.Port = *confPort
	}

	if "" != *confBaiduToken {
		Conf.BaiduToken = *confBaiduToken
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return tablePrefix + defaultTableName
	}

	logger.Debugf("configurations [%#v]", Conf)
}
