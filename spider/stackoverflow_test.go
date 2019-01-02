// 协慌网 https://routinepanic.com
// Copyright (C) 2018-2019, b3log.org

// 协慌网 https://routinepanic.com
// Copyright (C) 2018, b3log.org

package spider

import (
	"testing"
)

func TestParseQuestion(t *testing.T) {
	qna := StackOverflow.ParseQuestion("https://stackoverflow.com/questions/11227809/why-is-it-faster-to-process-a-sorted-array-than-an-unsorted-array")
	t.Log("question id: " + qna.Question.SourceID)
	for _, answer := range qna.Answers {
		t.Log("answer id: " + answer.SourceID)
	}
}

func TestParseQuestions(t *testing.T) {
	qnas := StackOverflow.ParseQuestions("https://stackoverflow.com/questions?sort=votes")
	for _, qna := range qnas {
		t.Log(qna.Question.TitleEnUS)
	}
}
