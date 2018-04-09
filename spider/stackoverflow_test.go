package spider

import (
	"testing"
)

func TestParseQuestion(t *testing.T) {
	question, answers := StackOverflow.ParseQuestion("https://stackoverflow.com/questions/11227809/why-is-it-faster-to-process-a-sorted-array-than-an-unsorted-array")
	t.Log(question.SourceID)
	for _, answer := range answers {
		t.Log(answer.SourceID)
	}
}

func TestParseQuestions(t *testing.T) {
	qnas := StackOverflow.ParseQuestions("https://stackoverflow.com/questions?sort=votes")
	for _, qna := range qnas {
		t.Log(qna.question.Title)
	}
}