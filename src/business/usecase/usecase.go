package usecase

import (
	log "github.com/mafazans/kumparan/lib/logger"

	"github.com/mafazans/kumparan/src/business/domain"
	"github.com/mafazans/kumparan/src/business/usecase/article"
)

type Usecase struct {
	Article article.UsecaseItf
}

func Init(logger log.Logger, dom *domain.Domain) *Usecase {
	return &Usecase{
		Article: article.InitArticleUsecase(
			logger,
			dom.Article,
		),
	}
}
