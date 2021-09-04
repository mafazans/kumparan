package restserver

import (
	log "github.com/mafazans/kumparan/lib/logger"

	"net/http"
	"sync"

	"github.com/mafazans/kumparan/src/business/usecase"

	"github.com/gorilla/mux"
)

var once = &sync.Once{}

type REST interface{}

type rest struct {
	logger log.Logger
	uc     *usecase.Usecase
	router *mux.Router
}

func Init(logger log.Logger, uc *usecase.Usecase, router *mux.Router) REST {
	var e *rest
	once.Do(func() {
		e = &rest{
			logger: logger,
			uc:     uc,
			router: router,
		}
		e.Serve()
	})
	return e
}

func (e *rest) Serve() {
	e.router.HandleFunc(`/`, e.home).Methods(`GET`)
	e.router.HandleFunc(`/article`, e.createArticle).Methods(`POST`)
	e.router.HandleFunc(`/article`, e.getArticle).Methods(`GET`)
	e.router.HandleFunc(`/article/{id}`, e.getArticleByID).Methods(`GET`)

	http.ListenAndServe(`:8000`, e.router)
}
