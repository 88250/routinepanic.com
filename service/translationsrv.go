package service

import "github.com/parnurzeal/gorequest"

// Translation service.
var Translation = &translationService{}

type translationService struct {
}

func (srv *translationService) Translate(text string) string {
	_, ret, errs := gorequest.New().Post("http://47.89.254.198:6868").SendMap(map[string]string{"text": text}).End()
	if nil != errs {
		logger.Errorf("translate failed: %s", errs)

		return ""
	}

	return ret
}
