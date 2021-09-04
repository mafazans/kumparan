package restserver

import "github.com/mafazans/kumparan/src/business/entity"

type HTTPArticleResp struct {
	Meta entity.Meta `json:"metadata"`
	Data ArticleData `json:"data"`
}

// TokenInfoData
type ArticleData struct {
	Article entity.Article `json:"article"`
}

type HTTPErrResp struct {
	Meta entity.Meta `json:"metadata"`
}

type HTTPGetArticleResp struct {
	Meta entity.Meta    `json:"metadata"`
	Data GetArticleData `json:"data"`
}

type GetArticleData struct {
	Articles []entity.Article `json:"articles"`
}
