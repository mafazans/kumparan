package domain

import (
	"database/sql"

	log "github.com/mafazans/kumparan/lib/logger"
	"github.com/mafazans/kumparan/src/business/domain/article"

	"github.com/go-redis/redis"
)

type Domain struct {
	Article article.DomainItf
}

func Init(
	logger log.Logger,
	sql sql.DB,
	redis redis.Client) *Domain {
	return &Domain{
		// Article Domain
		Article: article.InitArticleDomain(
			logger,
			sql,
			redis,
		),
	}
}
