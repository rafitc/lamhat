package repository

import (
	"backend/src/core"
	customErrors "backend/src/errors"
	"backend/src/model"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

var Sugar_logger = core.Sugar

func FindUserByEmail(ctx *gin.Context, email string, tx pgx.Tx) (model.User, error) {
	var user model.User
	const query string = "SELECT * from app.users where email_id = $1"
	row := tx.QueryRow(ctx, query, email)
	err := row.Scan(&user)
	if err != nil {
		if err == pgx.ErrNoRows {
			return user, customErrors.ErrUserNotFound
		}
		// Handle other errors
		Sugar_logger.Errorf("Error querying user: %v", err)
		return user, err
	}
	return user, nil
}
