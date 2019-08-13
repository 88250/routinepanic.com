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

package service

import (
	"os"
	"time"

	"github.com/b3log/gulu"
	"github.com/b3log/routinepanic.com/util"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // mysql
)

// Logger
var logger = gulu.Log.NewLogger(os.Stdout)

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
	db.DB().SetConnMaxLifetime(5 * time.Minute)
	//db.LogMode(true)
}

// DisconnectDB disconnects from the database.
func DisconnectDB() {
	if err := db.Close(); nil != err {
		logger.Errorf("Disconnect from database failed: " + err.Error())
	}
}
