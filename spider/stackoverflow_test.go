package spider

import (
	"testing"
)

func TestParseQuestion(t *testing.T) {
	q := StackOverflow.ParseQuestion("https://stackoverflow.com/questions/11227809/why-is-it-faster-to-process-a-sorted-array-than-an-unsorted-array")
	t.Log(q)
}
