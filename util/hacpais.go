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

package util

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/parnurzeal/gorequest"
)

// CommunityURL is the URL of HacPai community.
const CommunityURL = "https://ld246.com"

// HacPaiUserInfo returns HacPai community user info specified by the given access token.
func HacPaiUserInfo(accessToken string) (ret map[string]interface{}) {
	result := map[string]interface{}{}
	response, data, errors := gorequest.New().TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		Post(CommunityURL+"/user/ak").SendString("access_token="+accessToken).Timeout(7*time.Second).
		Set("User-Agent", "Pipe; +https://github.com/88250/routinepanic.com").EndStruct(&result)
	if nil != errors || http.StatusOK != response.StatusCode {
		logger.Errorf("get community user info failed: %+v, %s", errors, data)
		return nil
	}
	if 0 != result["code"].(float64) {
		return nil
	}
	return result["data"].(map[string]interface{})
}
