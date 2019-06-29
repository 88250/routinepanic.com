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
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Pagination represents pagination info.
type Pagination struct {
	CurrentPageNum  int    `json:"currentPageNum"`
	PageSize        int    `json:"pageSize"`
	PageCount       int    `json:"pageCount"`
	WindowSize      int    `json:"windowSize"`
	RecordCount     int    `json:"recordCount"`
	PageNums        []int  `json:"pageNums"`
	NextPageNum     int    `json:"nextPageNum"`
	PreviousPageNum int    `json:"previousPageNum"`
	FirstPageNum    int    `josn:"firstPageNum"`
	LastPageNum     int    `json:"lastPageNum"`
	PageURL         string `json:"pageURL"`
}

// GetPage returns paging parameter.
func GetPage(c *gin.Context) int {
	ret, _ := strconv.Atoi(c.Query("p"))
	if 1 > ret {
		ret = 1
	}

	return ret
}

// NewPagination creates a new pagination with the specified current page num and record count.
func NewPagination(currentPageNum, recordCount int) *Pagination {
	pageCount := int(math.Ceil(float64(recordCount) / float64(PageSize)))

	previousPageNum := currentPageNum - 1
	if 1 > previousPageNum {
		previousPageNum = 0
	}
	nextPageNum := currentPageNum + 1
	if nextPageNum > pageCount {
		nextPageNum = 0
	}

	pageNums := paginate(currentPageNum, pageCount, WindowSize)
	firstPageNum := 0
	lastPageNum := 0
	if 0 < len(pageNums) {
		firstPageNum = pageNums[0]
		lastPageNum = pageNums[len(pageNums)-1]
	}

	return &Pagination{
		CurrentPageNum:  currentPageNum,
		NextPageNum:     nextPageNum,
		PreviousPageNum: previousPageNum,
		PageSize:        PageSize,
		PageCount:       pageCount,
		WindowSize:      WindowSize,
		RecordCount:     recordCount,
		PageNums:        pageNums,
		FirstPageNum:    firstPageNum,
		LastPageNum:     lastPageNum,
	}
}

func paginate(currentPageNum, pageCount, windowSize int) []int {
	var ret []int

	if pageCount < windowSize {
		for i := 0; i < pageCount; i++ {
			ret = append(ret, i+1)
		}
	} else {
		first := currentPageNum + 1 - windowSize/2
		if first < 1 {
			first = 1
		}
		if first+windowSize > pageCount {
			first = pageCount - windowSize + 1
		}
		for i := 0; i < windowSize; i++ {
			ret = append(ret, first+i)
		}
	}

	return ret
}
