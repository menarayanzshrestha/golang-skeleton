package middlewares

import (
	"net/http"

	"github.com/menarayanzshrestha/trello/constants"
	"github.com/menarayanzshrestha/trello/infrastructure"
	"github.com/menarayanzshrestha/trello/utils"

	"github.com/gin-gonic/gin"
)

//DBTransactionMiddleware -> struct for transaction
type DBTransactionMiddleware struct {
	handler infrastructure.Router
	logger  infrastructure.Logger
	db      infrastructure.Database
}

//NewDBTransactionMiddleware -> new instance of transaction
func NewDBTransactionMiddleware(
	handler infrastructure.Router,
	logger infrastructure.Logger,
	db infrastructure.Database,
) DBTransactionMiddleware {
	return DBTransactionMiddleware{
		handler: handler,
		logger:  logger,
		db:      db,
	}
}

//Handle -> It setup the database transaction middleware
func (m DBTransactionMiddleware) Handle() gin.HandlerFunc {
	m.logger.Zap.Info("setting up database transaction middleware")

	return func(c *gin.Context) {
		txHandle := m.db.DB.Begin()
		m.logger.Zap.Info("beginning database transaction")

		defer func() {
			if r := recover(); r != nil {
				txHandle.Rollback()
			}
		}()

		c.Set(constants.DBTransaction, txHandle)
		c.Next()

		if utils.StatusInList(c.Writer.Status(), []int{http.StatusOK, http.StatusCreated}) {
			m.logger.Zap.Info("committing transactions")
			if err := txHandle.Commit().Error; err != nil {
				m.logger.Zap.Error("trx commit error: ", err)
			}
		} else {
			m.logger.Zap.Info("rolling back transaction due to status code: ", c.Writer.Status())
			txHandle.Rollback()
		}
	}
}
