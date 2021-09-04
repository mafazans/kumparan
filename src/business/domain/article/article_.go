package article

import (
	"database/sql"
	"time"

	"github.com/mafazans/kumparan/src/business/entity"

	x "github.com/mafazans/kumparan/lib/errors"

	"github.com/go-redis/redis"
)

func (a *article) setArticleDefaultValue(v entity.Article) (entity.Article, error) {
	now := time.Now().UTC()
	v.CreatedAt = sql.NullTime{Valid: true, Time: now}
	return v, nil
}

func (a *article) CreateArticle(v entity.Article) (entity.Article, error) {
	// set default value
	v, err := a.setArticleDefaultValue(v)
	if err != nil {
		a.logger.Debug(x.Wrap(err, "err"))
		return v, err
	}
	v, err = a.createSQLArticle(v)
	if err != nil {
		return v, err
	}
	//save to redis
	if v, err = a.upsertCacheArticle(v); err != nil {
		a.logger.Debug(x.Wrap(err, "err"))
		// log here
	}

	return v, nil
}

func (a *article) GetArticle(c entity.CacheControl, p entity.ArticleParam) ([]entity.Article, error) {
	// if must revalidate
	if c.MustRevalidate {
		result, err := a.getSQLArticle(p)
		if err != nil {
			a.logger.Debug(x.Wrap(err, "err"))
			return result, err
		}

		if _, err := a.upsertCacheArticleByQuery(p, result); err != nil {
			// log here, no need handle error
			a.logger.Debug(x.Wrap(err, "err"))
		}
		return result, nil
	}

	// get from redis
	result, err := a.getCacheArticleByQuery(p)
	if err == redis.Nil {
		result, err := a.getSQLArticle(p)
		if err != nil {
			a.logger.Debug(x.Wrap(err, "err"))
			return result, err
		}

		if _, err := a.upsertCacheArticleByQuery(p, result); err != nil {
			// log here
			a.logger.Debug(x.Wrap(err, "err"))
		}

		return result, nil

	} else if err != nil {
		a.logger.Debug(x.Wrap(err, "err"))
		// log here
		// fallback if there is redis error e.g. bad conn, etc.
		return a.getSQLArticle(p)
	}
	return result, nil
}

// GetArticleByID
func (a *article) GetArticleByID(c entity.CacheControl, vid int64) (entity.Article, error) {

	// if must revalidate
	if c.MustRevalidate {
		result, err := a.getSQLArticleByID(vid)
		if err != nil {
			return result, err
		}

		if _, err := a.upsertCacheArticle(result); err != nil {
			a.logger.Debug(x.Wrap(err, "err"))
		}

		return result, nil
	}

	result, err := a.getCacheArticleByID(vid)
	if err == redis.Nil {
		a.logger.Debug(x.Wrap(err, "err"))
		result, err := a.getSQLArticleByID(vid)
		if err != nil {
			return result, err
		}

		if _, err := a.upsertCacheArticle(result); err != nil {
			a.logger.Debug(x.Wrap(err, "err"))
		}

		return result, nil

	} else if err != nil {
		a.logger.Debug(x.Wrap(err, "err"))
		return a.getSQLArticleByID(vid)
	}

	return result, nil
}
