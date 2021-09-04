package sql

import (
	"database/sql"
	"time"

	log "github.com/mafazans/kumparan/lib/logger"

	x "github.com/mafazans/kumparan/lib/errors"

	_ "github.com/go-sql-driver/mysql"
)

func InitSQL(logger log.Logger) (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/kumparan_db?parseTime=true")
	if err != nil {
		return db, x.WrapWithCode(err, 500, "Failed to open connection to Mysql.")
	}

	if db == nil {
		return db, x.NewWithCode(500, "Sql failed to init")
	}

	for {
		logger.Info("Waiting for SQL Initialize")

		err = db.Ping()
		if err == nil {
			logger.Info("SQL Initialized")
			break
		}

		time.Sleep(1 * time.Second)
		continue
	}

	return db, nil
}
