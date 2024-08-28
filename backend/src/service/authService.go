package service

import (
	"backend/src/core"
	"backend/src/model"
)

var Sugar_logger = core.Sugar

func SingUpService(body model.SignupBody) model.Response {
	var result model.Response

	Sugar_logger.Info("Signup Service started with body", body.EmailId)

	// Check whether the user exist in DB or not
	
	return result
}
