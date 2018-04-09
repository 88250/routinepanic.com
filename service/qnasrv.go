package service

import (
	"github.com/b3log/routinepanic.com/spider"
	"github.com/b3log/routinepanic.com/model"
)

// QnA service.
var QnA = &qnaService{
}

type qnaService struct {
}

func (src *qnaService) Add(qnas []*spider.QnA) (err error) {
	tx := db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	for _, qna := range qnas {
		if err = db.Where("`source_id` = ? AND `source` = ?", qna.Question.SourceID, qna.Question.Source).
			Assign(model.Question{
			Title:     qna.Question.Title,
			Tags:      qna.Question.Tags,
			Content:   qna.Question.Content,
			SourceURL: qna.Question.SourceURL,
		}).FirstOrCreate(qna.Question).Error; nil != err {
			return
		}
		for _, answer := range qna.Answers {
			if err = db.Where("`question_id` = ? AND `source` = ?", qna.Question.ID, answer.Source).
				Assign(model.Answer{
				Content:   answer.Content,
				SourceURL: answer.SourceURL,
			}).FirstOrCreate(answer).Error; nil != err {
				return
			}
		}
	}

	return nil
}
