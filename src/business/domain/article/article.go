package article

import (
	"database/sql"

	log "github.com/mafazans/kumparan/lib/logger"
	"github.com/mafazans/kumparan/src/business/entity"

	"github.com/go-redis/redis"
)

type DomainItf interface {
	CreateArticle(v entity.Article) (entity.Article, error)
	GetArticle(c entity.CacheControl, p entity.ArticleParam) ([]entity.Article, error)
	GetArticleByID(c entity.CacheControl, vid int64) (entity.Article, error)
}

type article struct {
	logger log.Logger
	sql    sql.DB
	redis  redis.Client
}

func InitArticleDomain(logger log.Logger, sql sql.DB, redis redis.Client) DomainItf {
	return &article{
		logger: logger,
		sql:    sql,
		redis:  redis,
	}
}
