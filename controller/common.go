package controller

import "GoNexus/common/code"

type Response struct {
	StatusCode code.Code `json:"status_code"`
	StatusMsg  string    `json:"status_msg,omitempty"`
}

func (res *Response) CodeOf(code code.Code) Response {
	res.StatusCode = code
	res.StatusMsg = code.Msg()
	return *res
}

func (res *Response) Success() {
	res.CodeOf(code.SuccessCode)
}
