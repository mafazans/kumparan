package article

import (
	"github.com/mafazans/kumparan/lib/common"
	x "github.com/mafazans/kumparan/lib/errors"
	"github.com/mafazans/kumparan/src/business/entity"

	"github.com/jmoiron/sqlx"
)

func (a *article) createSQLArticle(article entity.Article) (entity.Article, error) {
	tx, err := a.sql.Begin()
	if err != nil {
		a.logger.Debug(x.Wrap(err, "err"))
		return article, x.WrapWithCode(err, x.CodeSQLTxBegin, "CreateArticle")
	}
	row, err := tx.Exec(insertArticleQuery,
		article.Author,
		article.Title,
		article.Body,
		article.CreatedAt,
	)
	if err != nil {
		a.logger.Debug(x.Wrap(err, "err"))
		return article, x.WrapWithCode(err, x.CodeSQLCreate, "createSQLArticle")
	}
	if err := tx.Commit(); err != nil {
		a.logger.Debug(x.Wrap(err, "err"))
		return article, x.WrapWithCode(err, x.CodeSQLTxCommit, "createSQLArticle")
	}

	article.ID, err = row.LastInsertId()
	if err != nil {
		a.logger.Debug(x.Wrap(err, "err"))
		return article, x.WrapWithCode(err, x.CodeSQLCannotRetrieveLastInsertID, "CreateArticle")
	}

	return article, nil
}

func (a *article) getSQLArticle(p entity.ArticleParam) ([]entity.Article, error) {
	results := []entity.Article{}

	p.Limit = common.ValidateLimit(p.Limit)
	p.Page = common.ValidatePage(p.Page)

	// New Query Builder
	qb := common.NewSQLClauseBuilder("param", "db", ``, p.Page, p.Limit).
		AliasPrefix("article", &p)

	// Build Dynamic Query Extension
	queryExt, _, args, err := qb.Build()
	if err != nil {
		a.logger.Debug(x.Wrap(err, "err"))
		return results, x.WrapWithCode(err, x.CodeSQLBuilder, "GetArticle")
	}

	// Build Get Article Query
	query, args, err := sqlx.In(ReadArticleN+queryExt, args...)
	if err != nil {
		a.logger.Debug(x.Wrap(err, "err"))
		return results, x.WrapWithCode(err, x.CodeSQLBuilder, "GetArticle")
	}

	// Get Article
	rows, err := a.sql.Query(query, args...)
	if err != nil {
		a.logger.Debug(x.Wrap(err, "err"))
		return results, x.WrapWithCode(err, x.CodeSQLRead, "GetArticle")
	}
	defer rows.Close()

	for rows.Next() {
		var result entity.Article
		if err := rows.Scan(
			&result.ID,
			&result.Author,
			&result.Title,
			&result.Body,
			&result.CreatedAt,
		); err != nil {
			a.logger.Debug(x.Wrap(err, "err"))
			return results, x.WrapWithCode(err, x.CodeSQLRowScan, "GetArticle")
		}
		results = append(results, result)
	}

	return results, nil
}

func (a *article) getSQLArticleByID(vid int64) (entity.Article, error) {
	// in order to return empty object if contains nothing we have to initialize the object with zero values
	result := entity.Article{}

	row := a.sql.QueryRow(ReadArticle1, vid)

	err := row.Scan(
		&result.ID,
		&result.Author,
		&result.Title,
		&result.Body,
		&result.CreatedAt,
	)
	if err != nil {
		return result, x.WrapWithCode(err, x.CodeSQLRowScan, `GetArticleByID`)
	}

	return result, nil
}
