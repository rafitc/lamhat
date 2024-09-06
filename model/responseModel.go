package model

type Response struct {
	Status   bool        `json:status`
	Code     int         `json:"code"`
	Data     interface{} `json:"data, omitempty"`
	ErrorMsg string      `json:errorMsg`
}
