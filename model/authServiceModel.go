package model

import "time"

type User struct {
	Id               int       `json:"id" binding:"required"`
	Email_id         string    `json:"email_id" binding:"required"`
	Is_email_valid   bool      `json:"is_email_valid" binding:"required"`
	First_name       *string   `json:"first_name" binding:"optional"`
	Last_name        *string   `json:"last_name" binding:"optional"`
	Auth_key_hash    *string   `json:"auth_key_hash" binding:"optional"`
	Otp              *string   `json:"otp" binding:"optional"`
	Otp_generated_at time.Time `json:"otp_generated_at" binding:"optional"`
	Is_user_active   bool      `json:"is_user_active" binding:"optional"`
	Created_at       time.Time `json:"created_at" binding:"optional"`
	Last_updated_at  time.Time `json:"last_updated_at" binding:"optional"`
}

type LoginUserDTO struct {
	Id               int       `json:"id" binding:"required"`
	Email_id         string    `json:"email_id" binding:"required"`
	Is_user_active   bool      `json:"is_user_active" binding:"optional"`
	Otp_generated_at time.Time `json:"otp_generated_at" binding:"optional"`
}
