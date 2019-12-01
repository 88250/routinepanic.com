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
	"encoding/json"
	"time"

	"github.com/88250/routinepanic.com/model"
	"github.com/88250/routinepanic.com/util"
)

// Review service.
var Review = &reviewService{}

type reviewService struct {
}

func (srv *reviewService) FilterPassed(revisions []*model.Revision) (ret []*model.Revision) {
	ret = []*model.Revision{}
	for _, revision := range revisions {
		count := 0
		db.Model(&model.Review{}).Where("`revision_id` = ? AND `status` = ?", revision.ID, model.ReviewStatusPassed).
			Count(&count)
		if 0 < count {
			ret = append(ret, revision)
		}
	}

	return
}

func (srv *reviewService) RejectReview(review *model.Review) (err error) {
	tx := db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	review.ReviewedAt = time.Now()
	review.Status = model.ReviewStatusRejected

	if err = tx.Model(review).Update(review).Error; nil != err {
		return
	}

	return nil
}

func (srv *reviewService) PassReview(review *model.Review, arg map[string]interface{}) (err error) {
	revision := QnA.GetRevision(review.RevisionID)

	tx := db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	review.ReviewedAt = time.Now()
	review.Status = model.ReviewStatusPassed

	if err = tx.Model(review).Update(review).Error; nil != err {
		return
	}

	data := map[string]interface{}{}
	if err = json.Unmarshal([]byte(revision.Data), &data); nil != err {
		return
	}

	data["content"] = arg["content"].(string)

	if model.DataTypeQuestion == revision.DataType {
		data["title"] = arg["title"].(string)
		question := QnA.GetQuestionByID(revision.DataID)
		question.ContentZhCN = data["content"].(string)
		question.TitleZhCN = data["title"].(string)
		if err = tx.Model(question).Update(question).Error; nil != err {
			return
		}
	} else {
		answer := QnA.GetAnswerByID(revision.DataID)
		answer.ContentZhCN = data["content"].(string)
		if err = tx.Model(answer).Update(answer).Error; nil != err {
			return
		}
	}

	revisionBytes, _ := json.Marshal(data)
	revision.Data = string(revisionBytes)
	if err = tx.Model(revision).Update(revision).Error; nil != err {
		return
	}

	return nil
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
