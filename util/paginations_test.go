// 协慌网 https://routinepanic.com
// Copyright (C) 2018-2019, b3log.org

package util

import (
	"testing"
)

func TestPaginate(t *testing.T) {
	pageNums := paginate(1, 15, 20)
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	for i, val := range pageNums {
		if val != expected[i] {
			t.Errorf("exptected is [%d] at index [%d], actual is [%d]", expected[i], i, val)
		}
	}

	pageNums = paginate(50, 15, 20)
	expected = []int{41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60}
	for i, val := range pageNums {
		if val != expected[i] {
			t.Errorf("exptected is [%d] at index [%d], actual is [%d]", expected[i], i, val)
		}
	}
}
