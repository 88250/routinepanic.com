package service

import "testing"

func TestTranslate(t *testing.T) {
	text := Translate("测试<b>翻译</b>功能")
	t.Log(text)
}
