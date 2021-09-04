package restserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	x "github.com/mafazans/kumparan/lib/errors"
	"github.com/mafazans/kumparan/src/business/entity"
)

func (e *rest) httpRespSuccess(w http.ResponseWriter, r *http.Request, statusCode int, resp interface{}) {

	meta := entity.Meta{
		Path:       r.URL.String(),
		StatusCode: statusCode,
		Status:     http.StatusText(statusCode),
		Message:    fmt.Sprintf("%s %s [%d] %s", r.Method, r.URL.RequestURI(), statusCode, http.StatusText(statusCode)),
		Error:      nil,
		Timestamp:  time.Now().Format(time.RFC3339),
	}

	var (
		raw []byte
		err error
	)

	switch data := resp.(type) {
	case entity.Article:
		articleResp := &HTTPArticleResp{
			Meta: meta,
			Data: ArticleData{
				Article: data,
			},
		}

		raw, err = json.Marshal(articleResp)

	case []entity.Article:
		articleResp := &HTTPGetArticleResp{
			Meta: meta,
			Data: GetArticleData{
				Articles: data,
			},
		}

		raw, err = json.Marshal(articleResp)

	default:
		e.httpRespError(w, r, x.NewWithCode(500, fmt.Sprintf("cannot cast type of %+v", data)))
		return
	}

	if err != nil {
		e.httpRespError(w, r, x.WrapWithCode(err, 500, "MarshalHTTPResp"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(raw)
}

func (e *rest) httpRespError(w http.ResponseWriter, r *http.Request, err error) {
	statusCode, displayError := x.Compile(x.COMMON, err, "EN")

	jsonErrResp := &HTTPErrResp{
		Meta: entity.Meta{
			Path:       r.URL.String(),
			StatusCode: statusCode,
			Status:     http.StatusText(statusCode),
			Message:    fmt.Sprintf("%s %s [%d] %s", r.Method, r.URL.RequestURI(), statusCode, http.StatusText(statusCode)),
			Error:      &displayError,
			Timestamp:  time.Now().Format(time.RFC3339),
		},
	}

	raw, err := json.Marshal(jsonErrResp)
	if err != nil {
		statusCode = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(raw)
}
