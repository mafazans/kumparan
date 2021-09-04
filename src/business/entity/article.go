package entity

import "database/sql"

type Article struct {
	ID        int64        `param:"id" json:"id"`
	Author    string       `param:"author" json:"author"`
	Title     string       `param:"title" json:"title"`
	Body      string       `param:"body" json:"body"`
	CreatedAt sql.NullTime `param:"created_at" json:"created_at"`
}

type CreateArticleParam struct {
	Author string `json:"author"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type ArticleParam struct {
	ID        int64        `param:"id" json:"id" db:"id"`
	Author    string       `param:"author" json:"author" db:"author"`
	Title     string       `param:"title" json:"title" db:"title"`
	Body      string       `param:"body" json:"body" db:"body"`
	CreatedAt sql.NullTime `param:"created_at" json:"created_at" db:"created_at"`
	SortBy    []string     `param:"sort_by" db:"sort_by" json:"sort_by"`
	Page      int64        `param:"page" db:"page" json:"page"`
	Limit     int64        `param:"limit" db:"limit" json:"limit"`
}
