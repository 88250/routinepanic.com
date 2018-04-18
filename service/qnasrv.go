// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package service

import (
	"strings"

	"github.com/b3log/routinepanic.com/model"
	"github.com/b3log/routinepanic.com/spider"
	"github.com/b3log/routinepanic.com/util"
	"github.com/jinzhu/gorm"
)

// QnA service.
var QnA = &qnaService{}

type qnaService struct {
}

func (srv *qnaService) GetAnswers(questionID uint64) (ret []*model.Answer) {
	if err := db.Model(&model.Answer{}).Where("`question_id` = ?", questionID).Find(&ret).Error; nil != err {
		logger.Errorf("get answers of question [id=%d] failed: %s", questionID, err)

		return
	}

	return
}

func (srv *qnaService) GetQuestionByPath(path string) (ret *model.Question) {
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
		Find(&ret).Error; nil != err {
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
	if err = tx.Where("`source_id` = ? AND `source` = ?", qna.Question.SourceID, qna.Question.Source).
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

	if err = tagArticle(tx, qna.Question); nil != err {
		return
	}

	for _, answer := range qna.Answers {
		answer.QuestionID = qna.Question.ID
		if err = tx.Where("`question_id` = ? AND `source_id` = ? AND `source` = ?", qna.Question.ID, answer.SourceID, answer.Source).
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

func (srv *qnaService) UpdateSourceAll(qnas []*spider.QnA) (err error) {
	tx := db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	for _, qna := range qnas {
		if err = srv.updateSource(tx, qna); nil != err {
			return
		}
	}

	return nil
}

func (srv *qnaService) UpdateSource(qna *spider.QnA) (err error) {
	tx := db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	if err = srv.updateSource(tx, qna); nil != err {
		return
	}

	return nil
}

func (srv *qnaService) updateSource(tx *gorm.DB, qna *spider.QnA) (err error) {
	if err = tx.Model(qna.Question).Where("`source_id` = ? AND `source` = ?", qna.Question.SourceID, qna.Question.Source).
		Update(qna.Question).Error; nil != err {
		return
	}
	for _, answer := range qna.Answers {
		if err = tx.Model(answer).Where("`source_id` = ? AND `source` = ?", answer.SourceID, answer.Source).
			Update(answer).Error; nil != err {
			return
		}
	}

	return nil
}

func tagArticle(tx *gorm.DB, question *model.Question) error {
	tags := strings.Split(question.Tags, ",")
	for _, tagTitle := range tags {
		tag := &model.Tag{}
		tx.Where("`title` = ?", tagTitle).First(tag)
		if "" == tag.Title {
			tag.Title = tagTitle
			tag.QuestionCount = 1
			if err := tx.Create(tag).Error; nil != err {
				return err
			}
		} else {
			tag.QuestionCount += 1
			if err := tx.Model(tag).Updates(tag).Error; nil != err {
				return err
			}
		}

		rel := &model.Correlation{ID1: question.ID, ID2: tag.ID, Type: model.CorrelationQuestionTag}
		if err := tx.Create(rel).Error; nil != err {
			return err
		}
	}

	return nil
}
