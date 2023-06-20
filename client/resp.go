package client

import "github.com/yulecd/pp-common/perrors"

const (
	SuccessCode = 1
)

type Resp struct {
	Code      int         `json:"code"`
	Data      interface{} `json:"data"`
	Message   string      `json:"message"`
	Timestamp int64       `json:"timestamp"`
}

func (r *Resp) IsSuccess() bool {
	return r.Code == SuccessCode
}

func (r *Resp) GetError() error {
	if r.IsSuccess() {
		return nil
	}
	return perrors.GenError(r.Code, r.Message)
}
