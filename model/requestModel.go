package model

type SignupBody struct {
	EmailId string `json:"emailId" binding:"required"`
}
