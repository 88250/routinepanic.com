// 协慌网 https://routinepanic.com
// Copyright (C) 2018-2019, b3log.org

package service

import (
	"os"
	"time"

	"github.com/b3log/routinepanic.com/log"
	"github.com/b3log/routinepanic.com/util"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // mysql
)

// Logger
var logger = log.NewLogger(os.Stdout)

// ZeroTime represents zero time.
var ZeroTime, _ = time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")

var db *gorm.DB

// ConnectDB connects to the database.
func ConnectDB() {
	var err error

	db, err = gorm.Open("mysql", util.Conf.MySQL)
	if nil != err {
		logger.Fatalf("opens database failed: " + err.Error())
	}

	if err = db.AutoMigrate(util.Models...).Error; nil != err {
		logger.Fatal("auto migrate tables failed: " + err.Error())
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(50)
	//db.LogMode(true)
}

// DisconnectDB disconnects from the database.
func DisconnectDB() {
	if err := db.Close(); nil != err {
		logger.Errorf("Disconnect from database failed: " + err.Error())
	}
}
