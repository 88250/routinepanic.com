// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package service

import (
	"github.com/b3log/routinepanic.com/model"
	"github.com/b3log/routinepanic.com/util"
)

// Review service.
var Review = &reviewService{}

type reviewService struct {
}

func (srv *reviewService) PassReview() {

}

func (srv *reviewService) GetReviewByID(id uint64) (ret *model.Review) {
	ret = &model.Review{}
	if err := db.Where("`id` = ?", id).First(ret).Error; nil != err {
		logger.Errorf("get review [id=%d] failed: "+err.Error(), id)

		return nil
	}

	return
}

func (srv *reviewService) GetReviews(status, page int) (ret []*model.Review, pagination *util.Pagination) {
	offset := (page - 1) * util.PageSize
	count := 0

	if err := db.Model(&model.Review{}).
		Where("`status` = ?", status).
		Order("`created_at` DESC").Count(&count).
		Offset(offset).Limit(util.PageSize).
		Find(&ret).Error; nil != err {
		logger.Errorf("get waiting reviews failed: " + err.Error())
	}

	pagination = util.NewPagination(page, count)

	return
}
