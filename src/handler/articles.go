package restserver

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	x "github.com/mafazans/kumparan/lib/errors"
	"github.com/mafazans/kumparan/src/business/entity"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

func (e *rest) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Helloooo")
}

func (e *rest) createArticle(w http.ResponseWriter, r *http.Request) {
	var p entity.CreateArticleParam

	reqBody, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(reqBody, &p)
	if err != nil {
		e.httpRespError(w, r, x.WrapWithCode(err, x.CodeHTTPUnmarshal, "Failed to unmarshall create article param"))
		return
	}

	result, err := e.uc.Article.CreateArticle(entity.Article{
		Author: p.Author,
		Title:  p.Title,
		Body:   p.Body,
	})
	if err != nil {
		e.httpRespError(w, r, err)
		return
	}

	e.httpRespSuccess(w, r, http.StatusCreated, result)
}

func (e *rest) getArticle(w http.ResponseWriter, r *http.Request) {
	var p entity.ArticleParam

	// Make custom decoder so it can handle % as a wildcard
	query := strings.ReplaceAll(r.URL.RawQuery, "%", "%25")

	u, err := url.ParseQuery(query)
	if err != nil {
		e.httpRespError(w, r, x.WrapWithCode(err, x.CodeHTTPBadRequest, "ParseQuery"))
		return
	}

	if err := schema.NewDecoder().Decode(&p, u); err != nil {
		e.httpRespError(w, r, x.WrapWithCode(err, x.CodeHTTPBadRequest, "DecodeQueryParam"))
		return
	}

	var cacheControl entity.CacheControl
	if r.Header.Get("Cache-Control") == "must-revalidate" {
		cacheControl.MustRevalidate = true
	}

	// fetch article
	result, err := e.uc.Article.GetArticle(cacheControl, p)
	if err != nil {
		e.httpRespError(w, r, err)
		return
	}

	e.httpRespSuccess(w, r, http.StatusOK, result)
}

func (e *rest) getArticleByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idParam, ok := vars["id"]
	if !ok {
		e.httpRespError(w, r, x.NewWithCode(x.CodeHTTPBadRequest, "DecodeQueryParam"))
		return
	}
	_, err := regexp.MatchString("^[0-9]+$", idParam)
	if err != nil {
		e.httpRespError(w, r, x.WrapWithCode(err, x.CodeHTTPBadRequest, "DecodeQueryParam"))
		return
	}

	vid, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		e.httpRespError(w, r, x.WrapWithCode(err, x.CodeHTTPBadRequest, "DecodeQueryParam"))
		return
	}

	var cacheControl entity.CacheControl
	if r.Header.Get("Cache-Control") == "must-revalidate" {
		cacheControl.MustRevalidate = true
	}

	// fetch article
	result, err := e.uc.Article.GetArticleByID(cacheControl, vid)
	if err != nil {
		e.httpRespError(w, r, err)
		return
	}

	e.httpRespSuccess(w, r, http.StatusOK, result)
}
