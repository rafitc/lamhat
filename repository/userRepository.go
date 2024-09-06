package repository

import (
	"lamhat/core"
	customErrors "lamhat/errors"
	"lamhat/model"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

var Sugar_logger = core.Sugar

func FindUserByEmail(ctx *gin.Context, email string, tx pgx.Tx) (model.LoginUserDTO, error) {
	var user model.LoginUserDTO
	const query string = "SELECT id, email_id, is_user_active, otp_generated_at FROM app.users WHERE email_id = $1;"

	row := tx.QueryRow(ctx, query, email)
	err := row.Scan(
		&user.Id,               // Destination for 'id'
		&user.Email_id,         // Destination for 'email_id'
		&user.Is_user_active,   // Destination for 'is_user_active'
		&user.Otp_generated_at, // Destination for 'otp_generated_at'
	)

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

// Insert into DB
// For now only insert Email ID, later user can update their profile

func CreateNewUser(ctx *gin.Context, email string, otp string, tx pgx.Tx) (model.User, error) {
	var user model.User
	const query string = "INSERT INTO app.users (email_id, otp, otp_generated_at) VALUES ($1, $2, CURRENT_TIMESTAMP) RETURNING id, email_id, is_email_valid, first_name, last_name, auth_key_hash, otp, otp_generated_at, is_user_active, created_at, last_updated_at;"

	// Use QueryRow to get the returned user row
	err := tx.QueryRow(ctx, query, email, otp).Scan(
		&user.Id,
		&user.Email_id,
		&user.Is_email_valid,
		&user.First_name,
		&user.Last_name,
		&user.Auth_key_hash,
		&user.Otp,
		&user.Otp_generated_at,
		&user.Is_user_active,
		&user.Created_at,
		&user.Last_updated_at,
	)

	if err != nil {
		return user, err
	}

	// Commit the transaction after successful insertion
	if err := tx.Commit(ctx); err != nil {
		return user, err
	}

	return user, nil
}

// For login, update the otp and otp_generated_at time for exisitng user
func UpdateLoginOtp(ctx *gin.Context, user_id int, otp string, tx pgx.Tx) (bool, error) {
	const query string = "UPDATE app.users SET otp = $1, otp_generated_at = CURRENT_TIMESTAMP WHERE id = $2"
	commandTag, err := tx.Exec(ctx, query, otp, user_id)

	if err != nil {
		return false, err
	}

	if commandTag.RowsAffected() != 1 {
		return false, customErrors.LoginOtpUpdationFailed
	}
	tx.Commit(ctx)
	return true, nil
}
