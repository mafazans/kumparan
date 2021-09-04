package article

const (
	insertArticleQuery = `INSERT INTO article(
		author,
		title,
		body,
		created_at
	) VALUES(?, ?, ?, ?);`

	ReadArticleN = `SELECT 
		id,
		author,
		title,
		body,
		created_at
	FROM article`

	ReadArticle1 = `SELECT 
		id,
		author,
		title,
		body,
		created_at
	FROM article
	WHERE id=?;`
)
