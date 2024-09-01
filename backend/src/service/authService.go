package service

import (
	"backend/src/core"
	customErrors "backend/src/errors"
	"backend/src/model"
	"backend/src/repository"
	"backend/src/utils"

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

		result.Status = false
		result.Data = nil
		result.Code = 500
		result.ErrorMsg = err.Error()
		return result
	}
	defer connection.Release()

	// Get transaction from connection and use it till the end.
	// If any err, do rollback else do commit
	tx, err := connection.Begin(ctx)
	if err != nil {
		Sugar_logger.Errorf("Error in DB connection %v", err.Error())
		defer tx.Rollback(ctx)

		result.Status = false
		result.Data = nil
		result.Code = 500
		result.ErrorMsg = err.Error()
		return result
	}
	// Check whether the user exist in DB or not
	user, err := repository.FindUserByEmail(ctx, body.EmailId, tx)
	if err == customErrors.ErrUserNotFound {
		// Handle the case where no user is found
		Sugar_logger.Warnf("No user found with email id %s", body.EmailId)
		// Now create a new user
		// Generate OTP
		Sugar_logger.Infof("Generating OTP for the user %v", body.EmailId)
		otp, err := utils.GenerateOtp()
		if err != nil {
			Sugar_logger.Errorf("Error while creating OTP for the user %s : %v", body.EmailId, err.Error())

			result.Status = false
			result.Data = nil
			result.Code = 500
			result.ErrorMsg = err.Error()
			return result
		}

		// Store otp, userEmail and otp generated time in table
		Sugar_logger.Infof("Registering user %v", body.EmailId)
		response, err := repository.CreateNewUser(ctx, body.EmailId, otp, tx)
		if err != nil {
			Sugar_logger.Errorf("Error while registering user %s : %v", body.EmailId, err.Error())

			result.Status = false
			result.Data = nil
			result.Code = 500
			result.ErrorMsg = err.Error()
			return result
		}

		result.Status = true
		result.Data = response
		result.Code = 200
		result.ErrorMsg = ""
		return result

	}
	// User already present, then generate a new otp and update the table
	// Generate OTP
	Sugar_logger.Infof("User already present with email %v", body.EmailId)
	otp, err := utils.GenerateOtp()
	if err != nil {
		Sugar_logger.Errorf("Error while creating OTP for the user %s : %v", body.EmailId, err.Error())

		result.Status = false
		result.Data = nil
		result.Code = 500
		result.ErrorMsg = err.Error()
		return result
	}

	// Update the otp
	status, err := repository.UpdateLoginOtp(ctx, user.Id, otp, tx)
	if err != nil {
		Sugar_logger.Errorf("Error while login for the user %s", user.Email_id)

		result.Status = false
		result.Data = nil
		result.Code = 500
		result.ErrorMsg = err.Error()
		return result
	}

	if status {
		Sugar_logger.Infof("Successfully sent OTP for login for user %s", user.Email_id)
	}
	// commit everything and then return a result
	tx.Commit(ctx)

	result.Status = true
	result.Data = user
	result.Code = 200
	result.ErrorMsg = ""
	return result
}
