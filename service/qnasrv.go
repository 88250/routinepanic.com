// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package service

import (
	"github.com/b3log/routinepanic.com/model"
	"github.com/b3log/routinepanic.com/spider"
	"github.com/b3log/routinepanic.com/util"
	"github.com/jinzhu/gorm"
)

// QnA service.
var QnA = &qnaService{}

type qnaService struct {
}

func (src *qnaService) GetQuestionByPath(path string) (ret *model.Question) {
	ret = &model.Question{}
	if err := db.Model(&model.Question{}).Where("`path` = ?", path).First(ret).Error; nil != err {
		logger.Errorf("get question [path=" + path + "] failed: " + err.Error())

		return nil
	}

	return
}

func (srv *qnaService) GetQuestions(page int) (ret []*model.Question, pagination *util.Pagination) {
	offset := (page - 1) * util.PageSize
	count := 0

	if err := db.Model(&model.Question{}).
		Select("`id`, `created_at`, `title_zh_cn`, `tags`, `path`").
		Where("`title_zh_cn` != '' AND `content_zh_cn` != ''").
		Order("`created_at` DESC").Count(&count).
		Offset(offset).Limit(util.PageSize).
		Find(&ret).Debug().Error; nil != err {
		logger.Errorf("get questions failed: " + err.Error())
	}

	pagination = util.NewPagination(page, count)

	return
}

func (srv *qnaService) AddAll(qnas []*spider.QnA) (err error) {
	tx := db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	for _, qna := range qnas {
		if err = srv.add(tx, qna); nil != err {
			return
		}
	}

	return nil
}

func (srv *qnaService) Add(qna *spider.QnA) (err error) {
	tx := db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	if err = srv.add(tx, qna); nil != err {
		return
	}

	return nil
}

func (srv *qnaService) add(tx *gorm.DB, qna *spider.QnA) (err error) {
	if err = db.Where("`source_id` = ? AND `source` = ?", qna.Question.SourceID, qna.Question.Source).
		Assign(model.Question{
			TitleEnUS:   qna.Question.TitleEnUS,
			TitleZhCN:   qna.Question.TitleZhCN,
			Tags:        qna.Question.Tags,
			ContentEnUS: qna.Question.ContentEnUS,
			ContentZhCN: qna.Question.ContentZhCN,
			SourceURL:   qna.Question.SourceURL,
		}).FirstOrCreate(qna.Question).Error; nil != err {
		return
	}
	for _, answer := range qna.Answers {
		answer.QuestionID = qna.Question.ID
		if err = db.Where("`question_id` = ? AND `source_id` = ? AND `source` = ?", qna.Question.ID, answer.SourceID, answer.Source).
			Assign(model.Answer{
				ContentEnUS: answer.ContentEnUS,
				ContentZhCN: answer.ContentZhCN,
				SourceURL:   answer.SourceURL,
			}).FirstOrCreate(answer).Error; nil != err {
			return
		}
	}

	return nil
}
