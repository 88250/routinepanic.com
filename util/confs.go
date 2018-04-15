// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

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

	"github.com/b3log/routinepanic.com/log"
	"github.com/b3log/routinepanic.com/model"
	"github.com/jinzhu/gorm"
)

// Logger
var logger = log.NewLogger(os.Stdout)

// Version of RP.
const Version = "1.0.0"

// Conf of RP.
var Conf *Configuration

// Models represents all models..
var Models = []interface{}{
	&model.Question{}, &model.Answer{},
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

	log.SetLevel(Conf.LogLevel)
	if "" != *confLogLevel {
		Conf.LogLevel = *confLogLevel
		log.SetLevel(*confLogLevel)
	}

	home, err := UserHome()
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

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return tablePrefix + defaultTableName
	}

	logger.Debugf("configurations [%#v]", Conf)
}
