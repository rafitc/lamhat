package model

type SignupBody struct {
	EmailId string `json:"emailId" binding:"required"`
}

type LoginBody struct {
	EmailId string `json:"emailId" binding:"required"`
	Otp     string `json:"otp" binding:"required"`
}
