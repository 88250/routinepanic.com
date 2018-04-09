package service

import "github.com/parnurzeal/gorequest"

func Translate(text string) string {
	_, ret, errs := gorequest.New().Post("http://localhost:6868").SendMap(map[string]string{"text": text}).End()
	if nil != errs {
		logger.Errorf("translate failed: %s", errs)

		return ""
	}

	return ret
}
