package article

import (
	"encoding/json"
	"fmt"
	"time"

	x "github.com/mafazans/kumparan/lib/errors"
	"github.com/mafazans/kumparan/src/business/entity"

	"github.com/go-redis/redis"
	"github.com/golang/snappy"
)

const (
	durationArticleExpiration time.Duration = 24 * time.Hour

	ArticleByID    = `article:id:%d`
	articleByQuery = `article:q:%s`
)

func (a *article) upsertCacheArticle(v entity.Article) (entity.Article, error) {
	// marshal to json
	rawJSON, err := json.Marshal(v)
	if err != nil {
		a.logger.Debug(x.Wrap(err, "err"))
		return v, x.WrapWithCode(err, x.CodeCacheMarshal, "UpsertArticle")
	}

	// snappy compression
	var encJSON []byte
	encJSON = snappy.Encode(encJSON, rawJSON)

	// build redis key
	key := fmt.Sprintf(ArticleByID, v.ID)

	if err := a.redis.Set(key, encJSON, durationArticleExpiration).Err(); err != nil {
		a.logger.Debug(x.Wrap(err, "err"))
		return v, x.WrapWithCode(err, x.CodeCacheSetSimpleKey, "UpsertArticle")
	}

	return v, nil
}

func (a *article) getCacheArticleByID(vid int64) (entity.Article, error) {
	result := entity.Article{}

	// build redis key
	key := fmt.Sprintf(ArticleByID, vid)

	// get key
	res, err := a.redis.Get(key).Bytes()
	if err == redis.Nil {
		a.logger.Debug(x.Wrap(err, "err"))
		return result, err
	} else if err != nil {
		a.logger.Debug(x.Wrap(err, "err"))
		return result, x.WrapWithCode(err, x.CodeCacheGetSimpleKey, "GetArticleByID")
	}

	// decode encoded json
	var decJSON []byte
	decJSON, err = snappy.Decode(decJSON, res)
	if err != nil {
		a.logger.Debug(x.Wrap(err, "err"))
		return result, x.WrapWithCode(err, x.CodeCacheDecode, "GetArticleByID")
	}

	// unmarshaling returned byte
	if err := json.Unmarshal(decJSON, &result); err != nil {
		a.logger.Debug(x.Wrap(err, "err"))
		return result, x.WrapWithCode(err, x.CodeCacheUnmarshal, "GetArticleByID")
	}

	return result, nil
}

func (a *article) upsertCacheArticleByQuery(p entity.ArticleParam, v []entity.Article) ([]entity.Article, error) {
	// serialize query param to string
	rawKey, err := json.Marshal(p)
	if err != nil {
		a.logger.Debug(x.Wrap(err, "err"))
		return v, x.WrapWithCode(err, x.CodeCacheMarshal, "ArticleParamByQuery")
	}

	// build key
	key := fmt.Sprintf(articleByQuery, string(rawKey))

	rawJSON, err := json.Marshal(v)
	if err != nil {
		a.logger.Debug(x.Wrap(err, "err"))
		return v, x.WrapWithCode(err, x.CodeCacheMarshal, "UpsertArticleByQuery")
	}

	// snappy compression on merchant
	var encJSON []byte
	encJSON = snappy.Encode(encJSON, rawJSON)

	// set key expiration
	if err := a.redis.Set(key, encJSON, durationArticleExpiration).Err(); err != nil {
		a.logger.Debug(x.Wrap(err, "err"))
		return v, x.WrapWithCode(err, x.CodeCacheSetSimpleKey, "UpsertArticleByQuery")
	}

	return v, nil
}

func (a *article) getCacheArticleByQuery(p entity.ArticleParam) ([]entity.Article, error) {
	var (
		results []entity.Article
	)

	// serialize query param to string
	rawKey, err := json.Marshal(p)
	if err != nil {
		a.logger.Debug(x.Wrap(err, "err"))
		return results, x.WrapWithCode(err, x.CodeCacheMarshal, "ArticleParam")
	}

	// build key
	key := fmt.Sprintf(articleByQuery, string(rawKey))

	// fetch merchants
	resultRaw, err := a.redis.Get(key).Bytes()
	if err == redis.Nil {
		a.logger.Debug(x.Wrap(err, "err"))
		return results, err
	} else if err != nil {
		a.logger.Debug(x.Wrap(err, "err"))
		return results, x.WrapWithCode(err, x.CodeCacheGetSimpleKey, "GetArticleByQuery")
	}

	// decode merchant (encoded json)
	var decJSON []byte
	decJSON, err = snappy.Decode(decJSON, resultRaw)
	if err != nil {
		a.logger.Debug(x.Wrap(err, "err"))
		return results, x.WrapWithCode(err, x.CodeCacheDecode, "GetArticleByQuery")
	}

	// unmarshaling returned byte
	if err := json.Unmarshal(decJSON, &results); err != nil {
		a.logger.Debug(x.Wrap(err, "err"))
		return results, x.WrapWithCode(err, x.CodeCacheUnmarshal, "GetArticleByQuery")
	}

	return results, nil
}
