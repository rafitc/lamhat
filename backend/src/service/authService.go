package service

import (
	"backend/src/core"
	customErrors "backend/src/errors"
	"backend/src/model"
	"backend/src/repository"

	"github.com/gin-gonic/gin"
)

var Sugar_logger = core.Sugar

func SingUpService(ctx *gin.Context, body model.SignupBody) model.Response {
	var result model.Response

	Sugar_logger.Info("Signup Service started with body", body.EmailId)

	// Acquire a connection from the pool
	connection, err := repository.ConObjOfDB.Acquire(ctx)
	if err != nil {
		Sugar_logger.Errorf("Error while acquiring connection from the database pool!! %v", err.Error())
		return result
	}
	defer connection.Release()

	// Get transaction from connection and use it till the end.
	// If any err, do rollback else do commit
	tx, err := connection.Begin(ctx)
	if err != nil {
		Sugar_logger.Errorf("Error in DB connection %v", err.Error())
		defer tx.Rollback(ctx)
	}
	// Check whether the user exist in DB or not
	response, err := repository.FindUserByEmail(ctx, body.EmailId, tx)
	if err == customErrors.ErrUserNotFound {
		// Handle the case where no user is found
		Sugar_logger.Warnf("No user found with email id %s", body.EmailId)
	}
	tx.Commit(ctx)
	Sugar_logger.Info("This is the result %v", response)
	return result
}
