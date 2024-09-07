package service

import (
	"fmt"
	constants "lamhat/const"
	"lamhat/core"
	customErrors "lamhat/errors"
	"lamhat/model"
	"lamhat/repository"
	"lamhat/utils"

	"github.com/gin-gonic/gin"
)

var Sugar_logger = core.Sugar

type LogInResp struct {
	AuthCode string `json:"authcode" binding:"required"`
}

func SignUpService(ctx *gin.Context, body model.SignupBody) model.Response {
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

		// Send OTP via email

		otp_subject := constants.SIGNUP_OTP_EMAIL_SUBJECT
		otp_body := fmt.Sprintf(constants.OTP_EMAIL_BODY, otp)
		err = utils.TriggerEmail(otp_subject, otp_body, body.EmailId)
		if err != nil {
			Sugar_logger.Errorf("Error while sending OTP email for user %s - %s", body.EmailId, err.Error())
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

func LoginService(ctx *gin.Context, body model.LoginBody) model.Response {
	var result model.Response
	Sugar_logger.Infof("Login Service started with body %s", body.EmailId)

	// Fetch OTP and otp creation date
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

		result.Status = false
		result.Data = nil
		result.Code = 500
		result.ErrorMsg = err.Error()
		return result
	}

	Sugar_logger.Debug("Comparing OTP")
	err = utils.ValidateOTP(user.Otp_generated_at, user.Otp, body.Otp)
	if err != nil {
		Sugar_logger.Errorf("Error while validating OTP of user %s", err.Error())

		result.Status = false
		result.Data = nil
		result.Code = 402
		result.ErrorMsg = err.Error()
		return result
	}
	Sugar_logger.Info("OTP validated")

	Sugar_logger.Info("Generating auth token")
	// Generate Auth Token
	token, err := utils.GenerateAuthToken(user.Id)
	if err != nil {
		result.Status = false
		result.Data = nil
		result.Code = 402
		result.ErrorMsg = err.Error()
		return result
	}
	// Geneate response body
	authToken := LogInResp{token}

	result.Status = true
	result.Data = authToken
	result.Code = 200
	result.ErrorMsg = ""
	return result
}
