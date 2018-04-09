package service

import "testing"

func TestTranslate(t *testing.T) {
	text := Translation.Translate("Why is it faster to process a sorted array than an unsorted array?")
	t.Log(text)
}
