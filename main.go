package main

import (
	redisx "github.com/mafazans/kumparan/lib/database/redis"
	sqlx "github.com/mafazans/kumparan/lib/database/sql"
	log "github.com/mafazans/kumparan/lib/logger"

	"github.com/mafazans/kumparan/src/business/domain"
	"github.com/mafazans/kumparan/src/business/usecase"
	resthandler "github.com/mafazans/kumparan/src/handler"

	"github.com/gorilla/mux"
)

var (
	logger log.Logger
	uc     *usecase.Usecase
	router *mux.Router
	dom    *domain.Domain
)

func main() {
	InitServices()
}

func InitServices() error {
	logger = log.Init(log.Options{
		Level:         "info",
		Formatter:     "json",
		Output:        "stdout",
		LogOutputPath: "/tmp/log/app.log",
	})

	sql1, err := sqlx.InitSQL(logger)
	if err != nil {
		return err
	}
	redis1 := redisx.InitRedis(logger)

	router = mux.NewRouter()
	// Business layer Initialization

	dom = domain.Init(logger, *sql1, *redis1)
	uc = usecase.Init(logger, dom)

	_ = resthandler.Init(logger, uc, router)

	defer sql1.Close()

	return nil
}
