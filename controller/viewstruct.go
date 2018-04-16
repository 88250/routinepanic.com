package controller

type question struct {
	Title string
	Tags []*tag
}

type tag struct {
	Title string
}
