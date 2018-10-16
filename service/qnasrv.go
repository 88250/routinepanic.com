// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package service

import (
	"encoding/json"
	"strings"

	"github.com/b3log/routinepanic.com/model"
	"github.com/b3log/routinepanic.com/spider"
	"github.com/b3log/routinepanic.com/util"
	"github.com/jinzhu/gorm"
	"github.com/xrash/smetrics"
)

// QnA service.
var QnA = &qnaService{}

type qnaService struct {
}

func (srv *qnaService) GetRevision(revisionID uint64) (ret *model.Revision) {
	ret = &model.Revision{}
	if err := db.Where("`id` = ?", revisionID).Find(&ret).Error; nil != err {
		return
	}

	return
}

func (srv *qnaService) QRevisions(question *model.Question) (ret []*model.Revision) {
	if err := db.Where("`data_id` = ? AND `data_type` = ?", question.ID, model.DataTypeQuestion).Find(&ret).Error; nil != err {
		return
	}

	return
}

func (srv *qnaService) ARevisions(answer *model.Answer) (ret []*model.Revision) {
	if err := db.Where("`data_id` = ? AND `data_type` = ?", answer.ID, model.DataTypeAnswer).Find(&ret).Error; nil != err {
		return
	}

	return
}

func (srv *qnaService) ContriAnswer(author *model.User, answer *model.Answer) (err error) {
	old := srv.GetAnswerByID(answer.ID)
	distance := smetrics.WagnerFischer(old.ContentZhCN, answer.ContentZhCN, 1, 1, 2)
	if 0 == distance {
		return
	}

	winkler := smetrics.JaroWinkler(old.ContentZhCN, answer.ContentZhCN, 0.7, 4)
	logger.Info(author.Name+" ", distance, " ", winkler)

	tx := db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	revisionData := map[string]interface{}{
		"content": answer.ContentZhCN,
	}
	revisionBytes, _ := json.Marshal(revisionData)
	revision := &model.Revision{
		DataType:    model.DataTypeAnswer,
		DataID:      answer.ID,
		Data:        string(revisionBytes),
		AuthorID:    author.ID,
		Distance:    distance,
		JaroWinkler: winkler,
	}
	if err = tx.Save(revision).Error; nil != err {
		return
	}

	review := &model.Review{
		RevisionID: revision.ID,
		Status:     model.ReviewStatusWaiting,
		ReviewedAt: ZeroTime,
	}
	if err = tx.Save(review).Error; nil != err {
		return
	}

	return nil
}

func (srv *qnaService) ContriQuestion(author *model.User, question *model.Question) (err error) {
	old := srv.GetQuestionByID(question.ID)
	titleDistance := smetrics.WagnerFischer(old.TitleZhCN, question.TitleZhCN, 1, 1, 2)
	contentDistance := smetrics.WagnerFischer(old.ContentZhCN, question.ContentZhCN, 1, 1, 2)
	if 0 == titleDistance && 0 == contentDistance {
		return
	}

	winkler := smetrics.JaroWinkler(old.ContentZhCN, question.ContentZhCN, 0.7, 4)
	logger.Info(author.Name+" ", titleDistance, " ", contentDistance, " ", winkler)

	tx := db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	revisionData := map[string]interface{}{
		"title":   question.TitleZhCN,
		"content": question.ContentZhCN,
	}
	revisionBytes, _ := json.Marshal(revisionData)
	revision := &model.Revision{
		DataType:    model.DataTypeQuestion,
		DataID:      question.ID,
		Data:        string(revisionBytes),
		AuthorID:    author.ID,
		Distance:    contentDistance,
		JaroWinkler: winkler,
	}
	if err = tx.Save(revision).Error; nil != err {
		return
	}

	review := &model.Review{
		RevisionID: revision.ID,
		Status:     model.ReviewStatusWaiting,
		ReviewedAt: ZeroTime,
	}

	if err = tx.Save(review).Error; nil != err {
		return
	}

	return nil
}

func (srv *qnaService) GetAnswerByID(id uint64) (ret *model.Answer) {
	ret = &model.Answer{}
	if err := db.Model(&model.Answer{}).Where("`id` = ?", id).First(ret).Error; nil != err {
		logger.Errorf("get answer [id=%d] failed: "+err.Error(), id)

		return nil
	}

	return
}

func (srv *qnaService) GetQuestionByID(id uint64) (ret *model.Question) {
	ret = &model.Question{}
	if err := db.Model(&model.Question{}).Where("`id` = ?", id).First(ret).Error; nil != err {
		logger.Errorf("get question [id=%d] failed: "+err.Error(), id)

		return nil
	}

	return
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

func (srv *qnaService) GetTagQuestions(tagID uint64, page int) (ret []*model.Question, pagination *util.Pagination) {
	var rels []*model.Correlation
	if err := db.Where("`id2` = ? AND `type` = ?", tagID, model.CorrelationQuestionTag).
		Find(&rels).Error; nil != err {
		return
	}

	var questionIDs []uint64
	for _, questionTagRel := range rels {
		questionIDs = append(questionIDs, questionTagRel.ID1)
	}

	offset := (page - 1) * util.PageSize
	count := 0

	if err := db.Model(&model.Question{}).
		Select("`id`, `created_at`, `title_zh_cn`, `tags`, `path`").
		Where("`id` IN (?)", questionIDs).
		Order("`created_at` DESC").Count(&count).
		Offset(offset).Limit(util.PageSize).
		Find(&ret).Error; nil != err {
		logger.Errorf("get tag questions failed: " + err.Error())
	}

	pagination = util.NewPagination(page, count)

	return
}

func (srv *qnaService) GetQuestions(page int) (ret []*model.Question, pagination *util.Pagination) {
	offset := (page - 1) * util.PageSize
	count := 0

	if err := db.Model(&model.Question{}).
		Select("`id`, `created_at`, `title_zh_cn`, `tags`, `path`").
		Where("`title_zh_cn` != '' AND `content_zh_cn` != ''").
		Order("`votes` DESC").
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

	if err = tagQuestion(tx, qna.Question); nil != err {
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

	logger.Info("added a QnA [" + qna.Question.TitleEnUS + ", " + qna.Question.TitleZhCN + "]")

	return nil
}

func (srv *qnaService) TagAll(questions []*model.Question) (err error) {
	tx := db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	for _, question := range questions {
		if err = removeTagQuestionRels(tx, question); nil != err {
			return
		}
		if err = tagQuestion(tx, question); nil != err {
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

func tagQuestion(tx *gorm.DB, question *model.Question) error {
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

func removeTagQuestionRels(tx *gorm.DB, question *model.Question) error {
	var rels []*model.Correlation
	if err := tx.Where("`id1` = ? AND `type` = ?",
		question.ID, model.CorrelationQuestionTag).Find(&rels).Error; nil != err {
		return err
	}
	for _, rel := range rels {
		tag := &model.Tag{}
		if err := tx.Where("`id` = ?", rel.ID2).First(tag).Error; nil != err {
			continue
		}
		tag.QuestionCount -= 1
		if err := tx.Save(tag).Error; nil != err {
			continue
		}
	}

	if err := tx.Where("`id1` = ? AND `type` = ?", question.ID, model.CorrelationQuestionTag).
		Delete(&model.Correlation{}).Error; nil != err {
		return err
	}

	return nil
}
