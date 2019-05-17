// 协慌网 https://routinepanic.com
// Copyright (C) 2018-present, b3log.org

// 协慌网 https://routinepanic.com
// Copyright (C) 2018-2019, b3log.org

package controller

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/b3log/routinepanic.com/model"
	"github.com/b3log/routinepanic.com/service"
	"github.com/b3log/routinepanic.com/util"
	"github.com/vinta/pangu"
)

type review struct {
	ID          uint64
	Title       string
	URL         string
	Contributor *contributor
	CreatedAt   time.Time
	Distance    int
	JaroWinkler float64
	DataType    int
	DataID      uint64
	Status      int
	Memo        string
	OldTitle    string
	OldContent  string
	Content     string
}

type question struct {
	ID           uint64
	Path         string
	Title        string
	Description  string
	Tags         []*tag
	Content      template.HTML
	SourceURL    string
	Contributors []*contributor
}

type tag struct {
	Title string
}

type contributor struct {
	Name           string
	Avatar         string
	ContriCount    int
	ContriDistance int
}

type answer struct {
	ID           uint64
	Content      template.HTML
	Contributors []*contributor
}

// 七牛图片处理样式，用于贡献者头像
const QiniuImgStyleAvatar = "imageView2/1/w/64/h/64/interlace/0/q/100"

func reviewsVos(rModels []*model.Review) (ret []*review) {
	for _, rModel := range rModels {
		r := reviewVo(rModel)
		ret = append(ret, r)
	}

	return
}

func reviewVo(rModel *model.Review) (ret *review) {
	ret = &review{
		ID:        rModel.ID,
		CreatedAt: rModel.CreatedAt,
		Memo:      rModel.Memo,
		Status:    rModel.Status,
	}

	revision := service.QnA.GetRevision(rModel.RevisionID)
	ret.Distance = revision.Distance
	ret.JaroWinkler = revision.JaroWinkler
	ret.DataType = revision.DataType
	ret.DataID = revision.DataID

	data := map[string]interface{}{}
	if err := json.Unmarshal([]byte(revision.Data), &data); nil == err {
		ret.Content = data["content"].(string)
		if model.DataTypeQuestion == revision.DataType {
			ret.Title = data["title"].(string)
		}
	} else {
		logger.Errorf("unmarshal json failed: " + err.Error())
	}

	if model.DataTypeQuestion == ret.DataType {
		qModel := service.QnA.GetQuestionByID(ret.DataID)
		ret.OldTitle = qModel.TitleZhCN
		ret.OldContent = qModel.ContentZhCN
		ret.URL = util.Conf.Server + "/questions/" + qModel.Path
	} else {
		aModel := service.QnA.GetAnswerByID(ret.DataID)
		qModel := service.QnA.GetQuestionByID(aModel.QuestionID)
		ret.OldTitle = qModel.TitleZhCN
		ret.OldContent = aModel.ContentZhCN
		ret.URL = fmt.Sprintf(util.Conf.Server+"/questions/"+qModel.Path+"/answers/%d", aModel.ID)
	}

	author := service.User.Get(revision.AuthorID)
	ret.Contributor = &contributor{
		Name:   author.Name,
		Avatar: author.Avatar + "?" + QiniuImgStyleAvatar,
	}

	return
}

func questionsVos(qModels []*model.Question) (ret []*question) {
	for _, qModel := range qModels {
		q := questionVo(qModel)
		ret = append(ret, q)
	}

	return
}

func questionVo(qModel *model.Question) (ret *question) {
	desc := ""
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(qModel.ContentZhCN))
	if nil != err {
		logger.Errorf("extract description failed: " + err.Error())
	} else {
		desc = doc.Find("p").First().Text()
	}

	content := qModel.ContentZhCN
	content = util.TuneHTML(content)
	title := pangu.SpacingText(qModel.TitleZhCN)

	ret = &question{
		ID:          qModel.ID,
		Path:        qModel.Path,
		Title:       title,
		Description: desc,
		Content:     template.HTML(content),
		SourceURL:   qModel.SourceURL,
	}
	tagStrs := strings.Split(qModel.Tags, ",")
	for _, tagTitle := range tagStrs {
		ret.Tags = append(ret.Tags, &tag{Title: tagTitle})
	}

	revisions := service.QnA.QRevisions(qModel)
	revisions = service.Review.FilterPassed(revisions)
	contributorMap := map[uint64]*contributor{}

	for _, revision := range revisions {
		contributorId := revision.AuthorID
		if val, ok := contributorMap[contributorId]; !ok {
			contributor := &contributor{}
			user := service.User.Get(revision.AuthorID)
			contributor.Name = user.Name
			contributor.Avatar = user.Avatar + "?" + QiniuImgStyleAvatar
			contributor.ContriCount = 1
			contributor.ContriDistance = revision.Distance
			contributorMap[contributorId] = contributor
		} else {
			val.ContriCount += 1
			val.ContriDistance += revision.Distance
		}
	}

	for _, val := range contributorMap {
		ret.Contributors = append(ret.Contributors, val)
	}

	return
}

func answersVos(aModels []*model.Answer) (ret []*answer) {
	for _, aModel := range aModels {
		a := answerVo(aModel)
		ret = append(ret, a)
	}

	return
}

func answerVo(aModel *model.Answer) (ret *answer) {
	content := aModel.ContentZhCN
	content = util.TuneHTML(content)

	ret = &answer{
		ID:      aModel.ID,
		Content: template.HTML(content),
	}

	revisions := service.QnA.ARevisions(aModel)
	revisions = service.Review.FilterPassed(revisions)
	contributorMap := map[uint64]*contributor{}

	for _, revision := range revisions {
		contributorId := revision.AuthorID
		if val, ok := contributorMap[contributorId]; !ok {
			contributor := &contributor{}
			user := service.User.Get(revision.AuthorID)
			contributor.Name = user.Name
			contributor.Avatar = user.Avatar + "?" + QiniuImgStyleAvatar
			contributor.ContriCount = 1
			contributor.ContriDistance = revision.Distance
			contributorMap[contributorId] = contributor
		} else {
			val.ContriCount += 1
			val.ContriDistance += revision.Distance
		}
	}

	for _, val := range contributorMap {
		ret.Contributors = append(ret.Contributors, val)
	}

	return
}
