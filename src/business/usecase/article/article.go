package article

import (
	log "github.com/mafazans/kumparan/lib/logger"

	art "github.com/mafazans/kumparan/src/business/domain/article"
	"github.com/mafazans/kumparan/src/business/entity"
)

type UsecaseItf interface {
	CreateArticle(entity.Article) (entity.Article, error)
	GetArticle(c entity.CacheControl, p entity.ArticleParam) ([]entity.Article, error)
	GetArticleByID(c entity.CacheControl, vid int64) (entity.Article, error)
}

type article struct {
	logger log.Logger
	art    art.DomainItf
}

type Options struct {
}

func InitArticleUsecase(logger log.Logger, a art.DomainItf) UsecaseItf {
	return &article{
		logger: logger,
		art:    a,
	}
}

func (a *article) CreateArticle(p entity.Article) (entity.Article, error) {
	return a.art.CreateArticle(p)
}

func (a *article) GetArticle(c entity.CacheControl, p entity.ArticleParam) ([]entity.Article, error) {
	return a.art.GetArticle(c, p)
}

func (a *article) GetArticleByID(c entity.CacheControl, vid int64) (entity.Article, error) {
	return a.art.GetArticleByID(c, vid)
}
